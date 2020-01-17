from invoke import task


@task
def make_new_episode(c):
    """
    Generate new `./content/posts/podcast-$(next episode number).md` file
    """
    c.run("./make_new_episode.sh")


@task
def make_new_prep(c):
    """
    Generate new `./content/posts/prep-$(next episode number).md` file
    """
    c.run("./make_new_prep.sh")
