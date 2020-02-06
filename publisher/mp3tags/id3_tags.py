"""
Utility functions for working with id3 tags in podcast mp3 files
"""
import os.path
from datetime import date, timedelta
from string import Template
from typing import Dict, List

from eyed3 import id3, mimetype

from mp3tags.posts_parser import Chapter


def print_album_meta(tag: id3.Tag) -> None:
    meta_info_tags = {
        "Album": tag.album or "none",
        "Title": tag.title,
        "Cover Image": "<binary>" if tag.images else "none",
        "Artists": tag.artist or "none",
        "Episode": tag.track_num[0],
        "Released": tag.release_date,
        "Genre": tag.genre.name,
    }

    [print(f"{key}: {value}") for key, value in meta_info_tags.items()]


def print_chapter_frame(frame: id3.frames.ChapterFrame):
    print(f"{timedelta(seconds=frame.times[0] / 1000)} - {timedelta(seconds=frame.times[1] / 1000)}\t{frame.title}")


def print_toc(tag: id3.Tag) -> None:
    if not tag.table_of_contents:
        return

    print("Table of contents:")
    for toc in tag.table_of_contents:
        for chap_id in toc.child_ids:
            print_chapter_frame(tag.chapters[chap_id])


def set_mp3_table_of_contests(tag: id3.Tag, chapters: List[Chapter]):
    """
    Write table of contents to podcast episode mp3 file using id3 chapter frames
    (http://id3.org/id3v2-chapters-1.0)
    This TOC should be readable by Apple Podcasts.
    """
    if not chapters:
        raise RuntimeError("no table of contents received")

    if tag.table_of_contents:
        raise RuntimeError("File already have table of contents.")

    # write chapters info to file id3 metadata
    tag.table_of_contents.set(
        "toc".encode("ascii"),
        toplevel=True,
        child_ids=[f"chapter#{i}".encode("ascii") for i in range(0, len(chapters))],
        description="Темы",
    )

    prev_end = 0
    for item in chapters:
        assert item.start >= prev_end, f'Chapters are sorted incorrectly at "{item.title}"'
        prev_end = item.end

        tag.chapters.set(item.element_id, times=(item.start, item.end))
        added_chapter = tag.chapters.get(item.element_id)
        added_chapter.title = item.title


def set_mp3_album_tags(data: Dict[str, str], tag: id3.Tag, episode_num: int):
    # set album title and cover image
    tag.album = data["album"]
    image_type = id3.frames.ImageFrame.FRONT_COVER
    image_file = os.path.join("/srv", data["cover"].lstrip("/"))
    image_mime = mimetype.guessMimetype(image_file)

    # set cover image
    with open(image_file, "rb") as f:
        tag.images.set(image_type, f.read(), image_mime, "")

    # set various meta info
    tag.artist = data["artist"]
    tag.track_num = (episode_num, episode_num)
    tag.title = Template(data["title"]).substitute(episode_num=episode_num)
    tag.release_date = str(date.today())
    tag.genre = "Podcast"
