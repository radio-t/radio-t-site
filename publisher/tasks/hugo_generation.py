from invoke import task


@task
def new_episode(c):
    """
    Generate new `./content/posts/podcast-$(next episode number).md` file
    """
    c.run("./scripts/make_new_episode.sh")


@task
def new_prep(c):
    """
    Generate new `./content/posts/prep-$(next episode number).md` file
    """
    c.run("./scripts/make_new_prep.sh")
