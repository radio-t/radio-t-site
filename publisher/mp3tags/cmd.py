#!/usr/bin/env python3
import inspect
import os
import re
import sys

import click
import yaml
from eyed3 import core, id3

import mp3tags
from mp3tags.id3_tags import print_album_meta, print_toc, set_mp3_album_tags, set_mp3_table_of_contests
from mp3tags.posts_parser import parse_table_of_contents_from_md

EPISODES_DIRECTORY = os.getenv("LOCATION", "/episodes")


config_directory = os.path.dirname(inspect.getfile(mp3tags))
with open(os.path.join(config_directory, "config.yaml")) as f:
    CONFIG = yaml.safe_load(f)


@click.group()
def cli():
    """Command line tool to show and set mp3 tags to episode mp3 file"""
    pass


@cli.command()
@click.argument("episode", type=int)
@click.option("--location", help="podcast files location")
def print_tags(episode, location=''):
    """
    Print title, album, artist, ToC and other mp3 tags (relevant for Radio-T) in podcast episode file
    """
    location = location or os.getenv('LOCATION', '')
    path = os.path.join(location, f'rt_podcast{episode}', f'rt_podcast{episode}.mp3')

    full_path = _get_episode_mp3_full_path(path)

    episode_file = core.load(full_path)

    tag = episode_file.tag
    if not isinstance(tag, id3.Tag):
        print("Error: only ID3 tags can be extracted currently.", file=sys.stderr)
        sys.exit(1)

    print_album_meta(tag)
    print_toc(tag)


@cli.command()
@click.argument("episode", type=int)
@click.option("--location", help="podcast files location")
@click.option("--dry", is_flag=True, help="flag to dry-run, command will not save changes to the mp3 file")
@click.option("--verbose", is_flag=True, help="flag to show verbose output")
def set_tags(episode, location='', dry=False, verbose=False):
    """
    Add title, album, artists tags, set album image, write table of contents to podcast episode mp3 file.

    The ToC should be readable by Apple Podcasts.
    """
    location = location or os.getenv('LOCATION', '')
    path = os.path.join(location, f'rt_podcast{episode}', f'rt_podcast{episode}.mp3')
    full_path = _get_episode_mp3_full_path(path)

    # check that hugo template for new episode page is already exists
    # so we can parse table of contents from there
    episode_page_path = f"/srv/hugo/content/posts/podcast-{episode}.md"
    if not os.path.exists(episode_page_path):
        print(
            "Error:", f'New episode page "{episode_page_path}" does not exists', file=sys.stderr,
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

        set_mp3_album_tags(CONFIG["tags"], tag, episode)

        print("Parsing episode articles from markdown template for the episode page in `/hugo/content/posts/`...")

        toc = parse_table_of_contents_from_md(
            episode_page_path, CONFIG["toc"]["first_mp3_chapter_name"], CONFIG["toc"]["max_episode_hours"]
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


def _get_episode_mp3_full_path(path):
    """
    Return full path to podcast episode mp3 file located in `$EPISODES_DIRECTORY`
    """
    if not os.path.exists(EPISODES_DIRECTORY):
        print(
            "Error:", f'Directory "{EPISODES_DIRECTORY}" does not exists', file=sys.stderr,
        )
        sys.exit(1)

    full_path = os.path.join(EPISODES_DIRECTORY, path)
    if not os.path.exists(full_path):
        print("Error:", f'File "{full_path}" does not exists', file=sys.stderr)
        sys.exit(1)

    return full_path


if __name__ == "__main__":
    cli()
