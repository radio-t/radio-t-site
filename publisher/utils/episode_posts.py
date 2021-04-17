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


def parse_table_of_contents_from_md(filename: str, first_chapter_name: str, episode_length_secs: int) -> List[Chapter]:
    """
    Parse table of contents for episode from Hugo template for this episode post
    """
    # parse episode post markdown
    with open(filename, encoding="utf-8") as f:
        article_lines = [line for line in f.readlines() if line.lstrip().startswith("-")]

    article_regexp = re.compile(r"\-\s+?\[(.+?)\].*?\*([\d:]+)\*")
    articles = []
    prev_start = 0
    for line in article_lines:
        # parse article line, represent start time as an offset in seconds
        match_obj = article_regexp.match(line)
        if not match_obj:
            continue

        article, offset_str = match_obj.groups()
        article = article.strip()
        article_start = (dt.strptime(offset_str, "%H:%M:%S") - dt.strptime("00:00:00", "%H:%M:%S")).seconds
        assert article_start > prev_start, f'Themes are sorted incorrectly at "{line.strip()}"'
        prev_start = article_start

        articles.append((article, article_start))

    # insert an initial chapter - without it Apple Podcasts will show first chapter starting at 00:00:00
    # regardless of it's actual timings
    articles.insert(0, (first_chapter_name, 0))

    result = []
    for index, article_meta in enumerate(articles):
        article, start = article_meta

        if index + 1 < len(articles):
            end = articles[index + 1][1]
        else:
            # set last chapter end at total duration of the mp3 file
            end = episode_length_secs

        result.append(Chapter(element_id=new_id(index), title=article, start=start * 1000, end=end * 1000))

    return result
