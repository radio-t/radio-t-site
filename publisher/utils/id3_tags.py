"""
Utility functions for working with id3 tags in podcast mp3 files
"""
from datetime import date, timedelta
from typing import List

from eyed3 import id3, mimetype

from .episode_posts import Chapter


def print_album_meta(tag: id3.Tag) -> None:
    meta_info_tags = {
        "Album": tag.album or "none",
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


def set_episode_toc(tag: id3.Tag, chapters: List[Chapter]) -> None:
    """
    Write table of contents to podcast episode mp3 file using id3 chapter frames
    (http://id3.org/id3v2-chapters-1.0)
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


def set_mp3_table_of_contests(tag: id3.Tag, chapters: List[Chapter]):
    if not chapters:
        raise RuntimeError("no table of contents received")

    if tag.table_of_contents:
        raise RuntimeError("File already have table of contents.")

    set_episode_toc(tag, chapters)


def set_mp3_album_tags(tag: id3.Tag, filename: str, episode_num: int):
    # set album title and cover image
    tag.album = "Радио-Т"
    image_type = id3.frames.ImageFrame.FRONT_COVER
    image_file = "/srv/hugo/static/images/covers/cover.png"
    image_mime = mimetype.guessMimetype(image_file)

    # set cover image
    with open(image_file, "rb") as f:
        tag.images.set(image_type, f.read(), image_mime, "")

    # set various meta info
    tag.artist = "Umputun, Bobuk, Gray, Ksenks"
    tag.track_num = (episode_num, episode_num)
    tag.title = f"Радио-Т {episode_num}"
    tag.release_date = str(date.today())
    tag.genre = "Podcast"
