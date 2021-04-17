import os
import re
import sys

from eyed3 import core, id3
from invoke import task

from utils.episode_posts import parse_table_of_contents_from_md
from utils.id3_tags import print_album_meta, print_toc, set_mp3_album_tags, set_mp3_table_of_contests

EPISODES_DIRECTORY = os.getenv("EPISODES_DIRECTORY", "/episodes/")


@task(
    help={"path": f'podcast mp3 path relative to "{EPISODES_DIRECTORY}" directory in container beforehand'},
    auto_shortflags=False,
)
def print_mp3_tags(c, path):
    """
    Print title, album, artist, ToC and other mp3 tags (relevant for Radio-T) from podcast episode file.
    """
    full_path = _get_episode_mp3_full_path(path)

    episode_file = core.load(full_path)

    tag = episode_file.tag
    if not isinstance(tag, id3.Tag):
        print("Error: only ID3 tags can be extracted currently.", file=sys.stderr)
        sys.exit(1)

    print_album_meta(tag)
    print("ID3 tag header version:", ".".join(map(str, tag.header.version)))
    print_toc(tag)


@task(
    optional=["dry", "verbose"],
    help={
        "path": f'podcast mp3 path relative to "{EPISODES_DIRECTORY}" directory in container beforehand',
        "dry": "dry-run, running command will not save changes to the mp3 file",
        "verbose": "flag to show verbose output",
    },
    auto_shortflags=False,
)
def set_mp3_tags(c, path, dry=False, verbose=False):
    """
    Add title, album, artists tags, set album image, write table of contents to podcast episode mp3 file.
    The ToC should be readable by Apple Podcasts.
    """
    full_path = _get_episode_mp3_full_path(path)

    # check that hugo template for new episode page is already exists
    # so we can parse table of contents from there
    episode_num = int(re.match(r".*rt_podcast(\d*)\.mp3", path).group(1))
    episode_page_path = f"/srv/hugo/content/posts/podcast-{episode_num}.md"
    if not os.path.exists(episode_page_path):
        print(
            "Error:",
            f'New episode page "{episode_page_path}" does not exists',
            file=sys.stderr,
        )
        sys.exit(1)

    # remove both ID3 v1.x and v2.x tags.
    remove_version = id3.ID3_ANY_VERSION
    id3.Tag.remove(full_path, remove_version)

    episode_file = core.load(full_path)
    # using ID3v2.3 tags, because using newer ID3v2.4 version leads to problems with Apple Podcasts and Telegram
    # (they will stop showing chapters with long titles at all, see https://github.com/radio-t/radio-t-site/issues/209)
    episode_file.initTag(version=id3.ID3_V2_3)

    tag = episode_file.tag
    episode_length_secs = int(episode_file.info.time_secs)  # eyed3 returns episode length in float

    try:
        print("Creating new album meta tags: title, cover, artists, etc...")

        set_mp3_album_tags(dict(c.tags), tag, episode_num)

        print("Parsing episode articles from markdown template for the episode page in `/hugo/content/posts/`...")

        toc = parse_table_of_contents_from_md(episode_page_path, c.toc.first_mp3_chapter_name, episode_length_secs)

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


def _get_episode_mp3_full_path(path):
    """
    Return full path to podcast episode mp3 file located in `$EPISODES_DIRECTORY`
    """
    if not os.path.exists(EPISODES_DIRECTORY):
        print(
            "Error:",
            f'Directory "{EPISODES_DIRECTORY}" does not exists',
            file=sys.stderr,
        )
        sys.exit(1)

    full_path = os.path.join(EPISODES_DIRECTORY, path)
    if not os.path.exists(full_path):
        print("Error:", f'File "{full_path}" does not exists', file=sys.stderr)
        sys.exit(1)

    return full_path
