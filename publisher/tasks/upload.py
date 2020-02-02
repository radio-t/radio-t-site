import os
import os.path
import re
import sys

import requests
from invoke import task
from requests.auth import HTTPBasicAuth

from .hugo_generation import get_last_podcast_number
from .mp3_tags import EPISODES_DIRECTORY, get_episode_mp3_full_path, set_mp3_tags


@task(
    optional=["dry", "verbose"],
    help={
        "path": f'podcast mp3 path relative to "{EPISODES_DIRECTORY}" directory in container beforehand',
        "dry": "dry-run, running command will not save changes to the mp3 file",
        "verbose": "flag to show verbose output",
    },
    auto_shortflags=False,
)
def upload_mp3(c, path, dry=False, verbose=False):
    """
    Upload episode mp3 file to radio-t.com, archives, and run ansible tasks
    All lines printed with `!notif` prefix will be send as local notification via makefile
    """
    episode_num = int(re.match(r".*rt_podcast(\d*)\.mp3", path).group(1))
    print(f"!notif: Radio-T detected #{episode_num}")

    set_mp3_tags(c, path, dry=dry, verbose=verbose)
    if verbose:
        print("\n")

    print("!notif: Radio-T tagged")

    filename = os.path.basepath(path)
    full_path = get_episode_mp3_full_path(path)
    ssh_args = "-v" if verbose else ""

    print("!notif: Uploading mp3 file")

    c.run(f"scp {ssh_args} {full_path} umputun@master.radio-t.com:/srv/master-node/var/media/{filename}", pty=True)

    print("!notif: Removing old media files")
    find_command = "find /srv/master-node/var/media -type f -mtime +60 -mtime -1200 -exec rm -vf '{}' ';'"
    c.run(
        f'ssh {ssh_args} umputun@master.radio-t.com "{find_command}"', pty=True,
    )

    print("!notif: Runing ansible tasks")
    docker_command = f"docker exec -i ansible /srv/deploy_radiot.sh {episode_num}"
    c.run(
        f'ssh {ssh_args} umputun@master.radio-t.com "{docker_command}"', pty=True,
    )

    print("!notif: Copying to hp-usrv (local) archives")
    c.run(
        f"scp -P 2222 {ssh_args} {full_path} umputun@192.168.1.24:/data/archive.rucast.net/radio-t/media/",
        pty=True,
    )

    print("!notif: Uploading to archive site")
    c.run(f"scp {ssh_args} {full_path} umputun@master.radio-t.com:/data/archive/radio-t/media/{filename}", pty=True)
    chmod_command = f"chmod 644 /data/archive/radio-t/media/{filename}"
    c.run(f'ssh {ssh_args} umputun@master.radio-t.com "{chmod_command}"', pty=True)

    print(f"!notif: All done for {filename}")


@task(
    optional=["verbose"],
    help={"verbose": "flag to show verbose output"},
    auto_shortflags=False,
)
def deploy(c, verbose=False):
    """
    Commit new episode page to git, post message to gitter-bot, and remove articles from news
    """
    auth = os.getenv("RT_NEWS_ADMIN", "").strip()
    if not auth or ":" not in auth:
        print("Error:", "RT_NEWS_ADMIN environment variable not set", file=sys.stderr)
        sys.exit(1)

    news_user, news_passwd = auth.split(":")

    last_episode_num = get_last_podcast_number(c.http.site_url, c.http.user_agent, c.http.timeout)
    if not last_episode_num:
        print("Error:", f"Last podcast episode page not found", file=sys.stderr)
        sys.exit(1)

    current_episode_num = last_episode_num + 1
    print(f"Current episode number: {current_episode_num}")

    root_path = "/srv/"
    os.chdir(root_path)

    print("Pushing new episode post")
    c.run("git pull", pty=True)
    c.run("git add .", pty=True)
    c.run(f'git commit -m "auto episode after {current_episode_num}" && git push', pty=True)

    ssh_args = "-v" if verbose else ""

    print("Running hugo generation")
    docker_command = "cd /srv/site.hugo && git pull && docker-compose run --rm hugo"
    c.run(
        f'ssh {ssh_args} umputun@master.radio-t.com "{docker_command}"', pty=True,
    )

    print("Calling gitter-bot")
    gitter_bot_command = f"docker exec -i gitter-bot /srv/gitter-rt-bot --super=Umputun --super=bobuk --super=ksenks --super=grayru --dbg --export-num={current_episode_num} --export-path=/srv/html"
    c.run(
        f'ssh {ssh_args} umputun@master.radio-t.com "{gitter_bot_command}"', pty=True,
    )

    print("Removing news articles")
    headers = {"User-Agent": c.http.user_agent}
    resp = requests.delete(
        "https://news.radio-t.com/api/v1/news/active/last/8",
        auth=HTTPBasicAuth(news_user, news_passwd),
        headers=headers,
        timeout=c.http.timeout,
    )

    if resp.status_code != 200:
        print(
            "Error:",
            f"https://news.radio-t.com responded with status code {resp.status_code} when removing articles",
            file=sys.stderr,
        )
        sys.exit(1)

    print("Done")
