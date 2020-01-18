from invoke import Collection

from .episode_info import print_last_rt_link, print_next_episode_number
from .hugo_generation import new_episode, new_prep
from .mp3_tags import print_mp3_tags, set_mp3_tags
from .upload import deploy, upload_mp3

ns = Collection()

tasks = [
    new_prep,
    new_episode,
    print_last_rt_link,
    print_next_episode_number,
    print_mp3_tags,
    set_mp3_tags,
    deploy,
    upload_mp3,
]

[ns.add_task(obj) for obj in tasks]
