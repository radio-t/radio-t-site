from invoke import Collection

from .episode_info import print_last_rt_link, print_next_episode_number
from .hugo_generation import make_new_episode, make_new_prep
from .mp3_tags import print_mp3_tags, set_mp3_tags

ns = Collection()

tasks = [
    make_new_prep,
    make_new_episode,
    print_last_rt_link,
    print_next_episode_number,
    print_mp3_tags,
    set_mp3_tags,
]

[ns.add_task(obj) for obj in tasks]
