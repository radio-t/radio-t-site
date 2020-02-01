import sys
from urllib.parse import urljoin

import requests
from bs4 import BeautifulSoup
from invoke import task


@task
def print_next_episode_number(c):
    """
    Print to stdout next podcast episode number parsed from https://radio-t.com/
    """
    headers = {"User-Agent": c.http.user_agent}
    resp = requests.get(c.http.site_url, headers=headers, timeout=c.http.timeout)
    soup = BeautifulSoup(resp.content, "html.parser")

    for item in soup.find_all("h2", class_="number-title"):
        if "podcast-" in item.a["href"]:
            last_podcast_num = int(item.a["href"].strip("/").rsplit("podcast-", 1)[-1])
            print(last_podcast_num + 1)
            return

    print("Error:", f"Last podcast episode page not found", file=sys.stderr)
    sys.exit(1)


@task
def print_last_rt_link(c):
    """
    Print to stdout last podcast episode link parsed from https://radio-t.com/
    """
    headers = {"User-Agent": c.http.user_agent}
    resp = requests.get(c.http.site_url, headers=headers, timeout=c.http.timeout)
    soup = BeautifulSoup(resp.content, "html.parser")

    for item in soup.find_all("h2", class_="number-title"):
        if "podcast-" in item.a["href"]:
            print(urljoin(c.http.site_url, item.a["href"]))
            return

    print("Error:", f"Link to last podcast episode page not found", file=sys.stderr)
    sys.exit(1)
