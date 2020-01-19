import os
import re
import sys

from eyed3 import core, id3
from invoke import task

from utils.episode_posts import parse_table_of_contents_from_md
from utils.id3_tags import print_album_meta, print_toc, set_mp3_album_tags, set_mp3_table_of_contests

EPISODES_DIRECTORY = os.getenv("EPISODES_DIRECTORY", "/episodes/")


@task(
    help={
        "filename": f'podcast mp3 file name. \
                      File must be placed into "{EPISODES_DIRECTORY}" directory in container beforehand',
    },
    auto_shortflags=False,
)
def print_mp3_tags(c, filename):
    """
    Print title, album, artist, ToC and other mp3 tags (relevant for Radio-T) from podcast episode file.
    """
    full_path = get_episode_mp3_full_path(filename)

    episode_file = core.load(full_path)

    tag = episode_file.tag
    if not isinstance(tag, id3.Tag):
        print("Error: only ID3 tags can be extracted currently.", file=sys.stderr)
        sys.exit(1)

    print_album_meta(tag)
    print_toc(tag)


@task(
    optional=["dry", "verbose"],
    help={
        "filename": f'podcast mp3 file name. \
                      File must be placed into "{EPISODES_DIRECTORY}" directory in container beforehand',
        "dry": "dry-run, running command will not save changes to the mp3 file",
        "verbose": "flag to show verbose output",
    },
    auto_shortflags=False,
)
def set_mp3_tags(c, filename, dry=False, verbose=False):
    """
    Add title, album, artists tags, set album image, write table of contents to podcast episode mp3 file.
    The ToC should be readable by Apple Podcasts.
    """
    full_path = get_episode_mp3_full_path(filename)

    # check that hugo template for new episode page is already exists
    # so we can parse table of contents from there
    episode_num = int(re.match(r".*rt_podcast(\d*)\.mp3", filename).group(1))
    episode_page_file_path = f"/srv/hugo/content/posts/podcast-{episode_num}.md"
    if not os.path.exists(episode_page_file_path):
        print(
            "Error:", f'New episode page "{episode_page_file_path}" does not exists', file=sys.stderr,
        )
        sys.exit(1)

    # remove both ID3 v1.x and v2.x tags.
    remove_version = id3.ID3_ANY_VERSION
    id3.Tag.remove(full_path, remove_version)

    episode_file = core.load(full_path)
    episode_file.initTag(version=id3.ID3_V2_4)

    tag = episode_file.tag

    try:
        print("Creating new album meta tags: title, cover, artists, etc...")

        set_mp3_album_tags(dict(c.tags), tag, filename, episode_num)

        print("Parsing episode themes from markdown template for the episode page in `/hugo/content/posts/`...")

        toc = parse_table_of_contents_from_md(
            episode_page_file_path, c.toc.first_mp3_chapter_name, c.toc.max_episode_hours
        )

        print("Generating table of contents...")

        set_mp3_table_of_contests(tag, toc)

    except Exception as exc:
        print("Error:", str(exc), file=sys.stderr)
        sys.exit(1)

    if not dry:
        tag.save(encoding="utf8")
        print("New mp3 tags are saved.")

    if verbose:
        print("\n")
        print_album_meta(tag)
        print_toc(tag)


def get_episode_mp3_full_path(filename):
    """
    Return full path to podcast episode mp3 file located in `$EPISODES_DIRECTORY`
    """
    if not os.path.exists(EPISODES_DIRECTORY):
        print(
            "Error:", f'Directory "{EPISODES_DIRECTORY}" does not exists', file=sys.stderr,
        )
        sys.exit(1)

    full_path = os.path.join(EPISODES_DIRECTORY, filename)
    if not os.path.exists(full_path):
        print("Error:", f'File "{full_path}" does not exists', file=sys.stderr)
        sys.exit(1)

    return full_path
