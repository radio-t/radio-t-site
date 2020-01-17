"""
Write table of contents to podcast episode mp3 file using id3 chapter frames http://id3.org/id3v2-chapters-1.0
This TOC should be readable by Apple Podcasts.
"""

import re
import sys
from copy import copy
from dataclasses import dataclass
from datetime import date
from datetime import datetime as dt
from typing import Any, Dict, List, Union

from dateutil.parser import parse as dateutil_parse
from eyed3 import core, id3, mimetype

new_id = lambda index: f"chapter#{index}".encode("ascii")


@dataclass
class Chapter:
    element_id: bytes
    title: str
    start: int
    end: int

    @staticmethod
    def print_frame(frame: id3.frames.ChapterFrame):
        print(f"-- from {frame.times[0]} to {frame.times[1]}: {frame.title}")


def print_toc(tag: id3.Tag) -> None:
    for toc in tag.table_of_contents:
        for chap_id in toc.child_ids:
            Chapter.print_frame(tag.chapters[chap_id])


def set_episode_toc(tag: id3.Tag, chapters: List[Chapter]) -> None:
    """
    Write table of contents to podcast episode mp3 file using id3 chapter frames http://id3.org/id3v2-chapters-1.0
    This TOC should be readable by Apple Podcasts.
    """
    # write chapters info to file id3 metadata
    tag.table_of_contents.set(
        "toc".encode("ascii"),
        toplevel=True,
        child_ids=[f"chapter#{i}".encode("ascii") for i in range(0, len(chapters))],
        description="Темы",
    )

    for item in chapters:
        tag.chapters.set(item.element_id, times=(item.start, item.end))
        added_chapter = tag.chapters.get(item.element_id)
        added_chapter.title = item.title


def parse_table_of_contents_from_md(filename: str) -> List[Chapter]:
    """
    Parse table of contents for episode from Hugo template for this episode post
    """

    # parse episode post markdown
    with open(filename, encoding="utf-8") as f:
        theme_lines = [line for line in f.readlines() if line.lstrip().startswith("-")]

    theme_regexp = re.compile(r"\-\s+?\[(.+?)\].*?\*([\d:]+)\*")
    themes = []
    for line in theme_lines:
        match_obj = theme_regexp.match(line)
        if not match_obj:
            continue

        theme, offset_str = match_obj.groups()
        theme_start = (dt.strptime(offset_str, "%H:%M:%S") - dt.strptime("00:00:00", "%H:%M:%S")).seconds
        themes.append((theme, theme_start))

    # insert an initial chapter - without it Apple Podcasts will show first chapter starting at 00:00:00
    # regardless of it's actual timings
    themes.insert(0, ("Вступление", 0))

    result = []
    for index, theme_meta in enumerate(themes):
        theme, start = theme_meta

        if len(themes) < index + 1:
            end = themes[index + 1][1]
        else:
            end = 4 * 60 * 60  # 4 hours

        result.append(Chapter(element_id=new_id(index), title=theme, start=start * 1000, end=end * 1000))

    return result


def set_mp3_table_of_contests(tag: id3.Tag, chapters: List[Chapter], verbose: bool):
    if not chapters:
        raise RuntimeError("no table of contents received")

    print("Setting table of contents")
    if tag.table_of_contents:
        print("File already have table of contents.", file=sys.stderr)
        if verbose:
            print_toc(tag)

        raise RuntimeError()

    set_episode_toc(tag, chapters)

    print("Table of contents set.")
    if verbose:
        print_toc(tag)


def set_mp3_album_tags(tag: id3.Tag, filename: str, episode_num: int, verbose: bool):
    # set album title and cover image
    tag.album = "Радио-Т"
    image_type = id3.frames.ImageFrame.FRONT_COVER
    image_file = "/srv/hugo/static/images/covers/cover.png"
    image_mime = mimetype.guessMimetype(image_file)

    print(f"Setting cover image {image_file}")
    with open(image_file, "rb") as f:
        tag.images.set(image_type, f.read(), image_mime, "")

    print("New cover image set")

    # set various meta info
    print(f"Setting artist and title tags")
    tag.artist = "Umputun, Bobuk, Gray, Ksenks"
    tag.track_num = (episode_num, episode_num)
    tag.title = f"Радио-Т {episode_num}"
    tag.release_date = str(date.today())
    tag.genre = "Podcast"
    print("Artist and title tags set")
