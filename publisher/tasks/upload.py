from invoke import task

from .mp3_tags import EPISODES_DIRECTORY, get_episode_mp3_full_path, set_mp3_tags


@task(
    optional=["overwrite", "verbose"],
    help={
        "filename": f'podcast mp3 file name. \
                      File must be placed into "{EPISODES_DIRECTORY}" directory in container beforehand',
        "dry": "dry-run, running command will not save changes to the mp3 file",
        "verbose": "flag to show verbose output",
    },
    auto_shortflags=False,
)
def upload_mp3(c, filename, dry=False, verbose=False):
    """
    Upload episode mp3 file to radio-t.com, archives, and run ansible tasks
    All lines printed with `!notif` prefix will be send as local notification via makefile
    """
    episode_num = int(filename.split("rt_podcast", 1)[-1].rsplit(".")[0])
    print(f"!notif: Radio-T detected #{episode_num}")

    set_mp3_tags(c, filename, dry=dry, verbose=verbose)
    if verbose:
        print("\n")

    print("!notif: Radio-T tagged")

    full_path = get_episode_mp3_full_path(filename)
    ssh_args = "-v" if verbose else ""

    print("!notif: Uploading mp3 file")

    c.run(f"scp {ssh_args} {full_path} umputun@master.radio-t.com:/srv/master-node/var/media/{filename}", pty=True)

    print("!notif: Removing old media files")
    c.run(
        f"ssh {ssh_args} umputun@master.radio-t.com"
        + " \"find /srv/master-node/var/media -type f -mtime +60 -mtime -1200 -exec rm -vf '{}' ';'\"",
        pty=True,
    )

    print("!notif: Runing ansible tasks")
    c.run(
        f'ssh {ssh_args} umputun@master.radio-t.com "docker exec -i ansible /srv/deploy_radiot.sh {episode_num}"',
        pty=True,
    )

    print("!notif: Copying to hp-usrv (local) archives")
    c.run(
        f"scp -P 2222 {ssh_args} {full_path} umputun@archives.umputun.com:/data/archive.rucast.net/radio-t/media/",
        pty=True,
    )

    print("!notif: Uploading to archive site")
    c.run(f"scp {ssh_args} {full_path} umputun@master.radio-t.com:/data/archive/radio-t/media/{filename}", pty=True)
    c.run(f'ssh {ssh_args} umputun@master.radio-t.com "chmod 644 /data/archive/radio-t/media/{filename}"', pty=True)

    print(f"!notif: All done for {filename}")


@task
def deploy(c):
    """
    Commit new episode page to git, post message to gitter-bot, and remove articles from news
    """
    c.run("./scripts/deploy.sh")
