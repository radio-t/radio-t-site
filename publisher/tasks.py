import os
import re
import sys
from datetime import date

from eyed3 import core, id3, mimetype
from invoke import Collection, task

from mp3_chapters import (parse_table_of_contents_from_md, set_mp3_album_tags,
                          set_mp3_table_of_contests)


@task
def make_new_episode(c):
    """
    Generate new `./content/posts/podcast-$(next episode number).md` file
    """
    c.run('./make_new_episode.sh')


@task
def make_new_prep(c):
    """
    Generate new `./content/posts/prep-$(next episode number).md` file
    """
    c.run('./make_new_prep.sh')


@task
def print_next_episode_number(c):
    """
    Print to stdout next podcast episode number parsed from https://radio-t.com/
    """
    result = c.run('curl https://radio-t.com/ | grep rt_podcast | head -n1', hide=True)
    num = int(result.stdout.split("rt_podcast")[1][:3])+1
    print(num)


@task
def print_last_rt_link(c):
    """
    Print to stdout last podcast episode link parsed from https://radio-t.com/
    FIXME: currently broken - markup from the site is parsed incorrectly
    """
    result = c.run('curl https://radio-t.com/ | grep podcast- | head -n1', hide=True)
    link = "https://radio-t.com" + result.stdout.split("\"")[3]
    print(link)


EPISODES_DIRECTORY = os.getenv('EPISODES_DIRECTORY', '/episodes/')


@task(optional=['overwrite', 'verbose'], 
      help={'filename': f'podcast mp3 file name. File must be placed into "{EPISODES_DIRECTORY}" directory beforehand',
            'verbose': 'flag to show verbose output'},
      auto_shortflags=False)
def set_mp3_tags(c, filename, verbose=False):
    """
    Add title, album, artists, album image, write table of contents to podcast episode mp3 file using id3 chapter frames http://id3.org/id3v2-chapters-1.0
    This TOC should be readable by Apple Podcasts.
    """
    if not os.path.exists(EPISODES_DIRECTORY):
        print('Error:', f'Directory "{EPISODES_DIRECTORY}" does not exists', file=sys.stderr)
        sys.exit(1)

    full_path = os.path.join(EPISODES_DIRECTORY, filename)
    if not os.path.exists(full_path):
        print('Error:', f'File "{full_path}" does not exists', file=sys.stderr)
        sys.exit(1)

    # remove both ID3 v1.x and v2.x tags.
    remove_version = id3.ID3_ANY_VERSION
    id3.Tag.remove(full_path, remove_version)
    
    episode_file = core.load(full_path)
    episode_file.initTag(version=id3.ID3_V2_4)

    tag = episode_file.tag
    if not isinstance(tag, id3.Tag):
        print('Error: only ID3 tags can be extracted currently.', file=sys.stderr)
        sys.exit(1)
    
    episode_num = int(re.match(r'.*rt_podcast(\d*)\.mp3', filename).group(1))

    try:
        # set album tags and cover image
        set_mp3_album_tags(tag, filename, episode_num, verbose)
        # set table of contents
        toc = parse_table_of_contents_from_md(f'/srv/hugo/content/posts/podcast-{episode_num}.md')
        set_mp3_table_of_contests(tag, toc, verbose)
    except Exception as exc:
        print('Error:', str(exc), file=sys.stderr)
        sys.exit(1)
    else:
        tag.save(encoding='utf8')
        print('New mp3 tags are set.')
