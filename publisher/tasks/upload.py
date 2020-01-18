from invoke import task


@task
def upload_mp3(c):
    """
    Upload episode mp3 file to radio-t.com, archives, and run ansible tasks 
    """
    c.run("./scripts/upload_mp3.sh")


@task
def deploy(c):
    """
    Commit new episode page to git, post message to gitter-bot, and remove articles from news
    """
    c.run("./scripts/deploy.sh")
