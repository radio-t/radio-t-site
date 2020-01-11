import os
from invoke import task, Collection


@task
def make_new_episode(c):
    """
    Generate new `./content/posts/podcast-$(next episode number).md` file
    """
    c.run('./make_new_episode.sh')


@task
def make_new_prep(c):
    """
    Generate new `./content/posts/prep-$(next episode number).md` file
    """
    c.run('./make_new_prep.sh')


@task
def print_next_episode_number(c):
    """
    Print to stdout next podcast episode number parsed from https://radio-t.com/
    """
    result = c.run('curl https://radio-t.com/ | grep rt_podcast | head -n1', hide=True)
    num = int(result.stdout.split("rt_podcast")[1][:3])+1
    print(num)


@task
def print_last_rt_link(c):
    """
    Print to stdout last podcast episode link parsed from https://radio-t.com/
    FIXME: currently broken - markup from the site is parsed incorrectly
    """
    result = c.run('curl https://radio-t.com/ | grep podcast- | head -n1', hide=True)
    link = "https://radio-t.com" + result.stdout.split("\"")[3]
    print(link)