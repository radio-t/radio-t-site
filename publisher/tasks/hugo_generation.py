import os.path
import sys
from datetime import datetime as dt
from string import Template

import requests
from invoke import task

from .episode_info import get_last_episode_link, get_last_podcast_number

TEMPLATES_DIR = os.path.join(os.path.abspath(os.path.dirname(__file__)), "..", "templates")


@task
def new_episode(c):
    """
    Generate new `./content/posts/podcast-$(next episode number).md` file
    """
    # get new episode number from https://radio-t.com/
    last_episode_num = get_last_podcast_number(c.http.site_url, c.http.user_agent, c.http.timeout)
    if not last_episode_num:
        print("Error:", f"Last podcast episode page not found", file=sys.stderr)
        sys.exit(1)

    next_episode_num = last_episode_num + 1
    print(f"New post number: {next_episode_num}")

    # get template for new episode post
    new_file_path_relative = Template(c.hugo.episode_post).substitute(episode_num=next_episode_num).lstrip("/")
    new_file_path = os.path.join("/srv", new_file_path_relative)
    with open(os.path.join(TEMPLATES_DIR, "new_episode_post.tmpl"), "r") as f:
        new_episode_template = Template(f.read())

    # get themes from https://news.radio-t.com API
    headers = {"User-Agent": c.http.user_agent}
    news_api_resp = requests.get(c.http.themes_url, headers=headers, timeout=c.http.timeout)

    # write new episode post to hugo posts directory
    with open(new_file_path, "w", encoding="utf-8") as f:
        f.write(
            new_episode_template.substitute(
                episode_num=next_episode_num,
                timestamp=dt.now().strftime("%Y-%m-%dT%H:%M:%S"),
                themes=news_api_resp.text.strip(),
            )
        )

    print("New episode post generated. File:")
    print(new_file_path_relative)


@task
def new_prep(c):
    """
    Generate new `./content/posts/prep-$(next episode number).md` file
    """
    # get new episode number from https://radio-t.com/
    last_episode_num = get_last_podcast_number(c.http.site_url, c.http.user_agent, c.http.timeout)
    if not last_episode_num:
        print("Error:", f"Last podcast episode page not found", file=sys.stderr)
        sys.exit(1)

    next_episode_num = last_episode_num + 1
    print(f"New post number: {next_episode_num}")

    # get template for new prep post
    new_file_path_relative = Template(c.hugo.prep_post).substitute(episode_num=next_episode_num).lstrip("/")
    new_file_path = os.path.join("/srv", new_file_path_relative)
    with open(os.path.join(TEMPLATES_DIR, "new_prep_post.tmpl"), "r") as f:
        new_prep_template = Template(f.read())

    # write new prep post to hugo posts directory
    with open(new_file_path, "w", encoding="utf-8") as f:
        f.write(
            new_prep_template.substitute(
                episode_num=next_episode_num,
                timestamp=dt.now().strftime("%Y-%m-%dT%H:%M:%S"),
            )
        )

    print("Next episode prep generated. File:")
    print(new_file_path_relative)