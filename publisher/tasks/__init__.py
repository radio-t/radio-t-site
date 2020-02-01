from invoke import Collection

from .episode_info import print_last_rt_link, print_next_episode_number
from .mp3_tags import print_mp3_tags, set_mp3_tags

tasks = [
    print_last_rt_link,
    print_next_episode_number,
    print_mp3_tags,
    set_mp3_tags,
]

ns = Collection()
[ns.add_task(obj) for obj in tasks]
