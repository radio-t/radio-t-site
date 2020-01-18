import sys
from typing import Union
from urllib.parse import urljoin

import requests
from bs4 import BeautifulSoup
from invoke import task


@task
def print_next_episode_number(c):
    """
    Print to stdout next podcast episode number parsed from https://radio-t.com/
    """
    last_episode_num = get_last_podcast_number(c.http.site_url, c.http.user_agent, c.http.timeout)
    if not last_episode_num:
        print("Error:", f"Last podcast episode page not found", file=sys.stderr)
        sys.exit(1)

    print(last_episode_num + 1)


def get_last_podcast_number(site_url: str, user_agent: str, timeout: int = 30) -> Union[int, None]:
    headers = {"User-Agent": user_agent}
    resp = requests.get(site_url, headers=headers, timeout=timeout)
    soup = BeautifulSoup(resp.content, "html.parser")

    for item in soup.find_all("h2", class_="number-title"):
        if "podcast-" in item.a["href"]:
            last_podcast_num = int(item.a["href"].strip("/").rsplit("podcast-", 1)[-1])
            return last_podcast_num


@task
def print_last_rt_link(c):
    """
    Print to stdout last podcast episode link parsed from https://radio-t.com/
    """
    last_podcast_page_link = get_last_episode_link(c.http.site_url, c.http.user_agent, c.http.timeout)
    if not last_podcast_page_link:
        print("Error:", f"Link to last podcast episode page not found", file=sys.stderr)
        sys.exit(1)

    print(last_podcast_page_link)


def get_last_episode_link(site_url: str, user_agent: str, timeout: int = 30) -> Union[str, None]:
    headers = {"User-Agent": user_agent}
    resp = requests.get(site_url, headers=headers, timeout=timeout)
    soup = BeautifulSoup(resp.content, "html.parser")

    for item in soup.find_all("h2", class_="number-title"):
        if "podcast-" in item.a["href"]:
            return urljoin(site_url, item.a["href"])
