"""
Utility functions for working with episode posts in `/hugo/content/posts/` directory
"""
import re
from dataclasses import dataclass
from datetime import datetime as dt
from typing import List

new_id = lambda index: f"chapter#{index}".encode("ascii")


@dataclass
class Chapter:
    element_id: bytes
    title: str
    start: int  # start of chapter offset, seconds
    end: int  # end of chapter offset, seconds


def parse_table_of_contents_from_md(filename: str, first_chapter_name: str, max_episode_hours: int) -> List[Chapter]:
    """
    Parse table of contents for episode from Hugo template for this episode post
    """
    # parse episode post markdown
    with open(filename, encoding="utf-8") as f:
        theme_lines = [line for line in f.readlines() if line.lstrip().startswith("-")]

    theme_regexp = re.compile(r"\-\s+?\[(.+?)\].*?\*([\d:]+)\*")
    themes = []
    prev_start = 0
    for line in theme_lines:
        # parse theme line, represent start time as an offset in seconds
        match_obj = theme_regexp.match(line)
        if not match_obj:
            continue

        theme, offset_str = match_obj.groups()
        theme = theme.strip()
        theme_start = (dt.strptime(offset_str, "%H:%M:%S") - dt.strptime("00:00:00", "%H:%M:%S")).seconds
        assert theme_start > prev_start, f'Themes are sorted incorrectly at "{line.strip()}"'
        prev_start = theme_start

        themes.append((theme, theme_start))

    # insert an initial chapter - without it Apple Podcasts will show first chapter starting at 00:00:00
    # regardless of it's actual timings
    themes.insert(0, (first_chapter_name, 0))

    result = []
    for index, theme_meta in enumerate(themes):
        theme, start = theme_meta

        if index + 1 < len(themes):
            end = themes[index + 1][1]
        else:
            # set last chapter end at estimated end of the episode
            end = max_episode_hours * 60 * 60

        result.append(Chapter(element_id=new_id(index), title=theme, start=start * 1000, end=end * 1000))

    return result
