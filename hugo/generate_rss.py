#!/usr/bin/env python3

"""
Скрипт для генерации rss-файлов

pip install pytoml mistune
"""

import glob
import subprocess

import mistune
import pytoml as toml
from datetime import datetime

POSTS_DIR = './content/posts'
SAVE_TO = '/srv/hugo/public'
# SAVE_TO = '/tmp'
DATA_RSS = './data/rss'
FEEDS = [
    {'name': 'podcast', 'title': 'Радио-Т',
     'image': 'https://radio-t.com/images/cover.jpg', 'count': 20, 'size': True},
    {'name': 'podcast-archives', 'title': 'Радио-Т Архивы',
     'image': 'https://radio-t.com/images/cover_rt_big_archive.png', 'count': 1000, 'size': False},
    {'name': 'podcast-archives-short', 'title': 'Радио-Т Архивы',
     'image': 'https://radio-t.com/images/cover_rt_big_archive.png', 'count': 25, 'size': False},
]


def parse_file(name, source):
    data, config_lines, config_attr = list(), list(), 0

    for line in source:
        if line == '+++':
            config_attr += 1
        elif config_attr == 1:
            config_lines.append(line)
        else:
            data.append(line)

    toml_data = '\n'.join(config_lines)
    conf = toml.loads(toml_data)
    date = datetime.strptime(conf['date'], "%Y-%m-%dT%H:%M:%S")
    url = 'p/{}/{}/'.format(date.strftime('%Y/%m/%d'), name)

    return {'created_at': date, 'url': url, 'config': conf, 'data': '\n'.join(data)}


def get_mp3_size(mp3file):
    size = subprocess.check_output(
        "curl -sI http://archive.rucast.net/radio-t/media/" + mp3file + " | grep Content-Length | awk '{print $2}'",
        shell=True).decode("utf-8")
    size = size.replace("\r\n", "").replace("\n", "")
    print(mp3file, size)
    return size


def run():
    print("generate rss")
    renderer = mistune.Renderer(escape=False)
    markdown = mistune.Markdown(renderer=renderer)

    # загружаем настройки
    with open('config.toml', encoding='utf-8') as f:
        mconfig = toml.load(f)

    # получаем все файлы
    posts = list()
    for post_file in glob.glob(POSTS_DIR + '/*.md'):
        with open(post_file, encoding='utf-8') as h:
            name = post_file.replace(POSTS_DIR, '').replace('.md', '').replace('\\', '')
            post = parse_file(name, h.read().splitlines())
            # пропускаем посты, которые не являются подкастами
            if 'categories' not in post['config']:
                continue
            if 'podcast' not in post['config']['categories']:
                continue
            posts.append(post)

    # сотируем по дате и получаем первые `COUNT` постов
    posts.sort(key=lambda x: x['created_at'], reverse=True)
    # posts = posts[:COUNT + 1]

    # генерируем каждый фид
    for feed in FEEDS:
        # шапка
        with open(DATA_RSS + '/head.xml', encoding='utf-8') as f:
            head = f.read()
        head = head.format(title=feed['title'], url=mconfig['baseurl'],
                           subtitle=mconfig['params']['subtitle'], description=mconfig['params']['longDescription'],
                           image=feed['image'])

        # ноги
        with open(DATA_RSS + '/foot.xml', encoding='utf-8') as f:
            foot = f.read()

        # генерация постов
        feed_posts = list()
        with open(DATA_RSS + '/{}.xml'.format(feed['name']), encoding='utf-8') as f:
            body = f.read()
            for post in posts:
                if len(feed_posts) > feed['count']:
                    break

                def attr(x):
                    return post['config'][x] if x in post['config'] else ''

                date = post['created_at'].strftime('%a, %d %b %Y %H:%M:%S EST')

                fsize = ""
                if feed['count'] < 30 and feed['size'] is True:
                    fsize = get_mp3_size(attr('filename') + ".mp3")

                item = body.format(title=post['config']['title'], content=markdown(post['data']),
                                   filename=attr('filename'),
                                   filesize=fsize,
                                   fixed_url='{}/{}'.format(mconfig['baseurl'], post['url']).replace("//p", "/p"),
                                   date=date, image=attr('image'), url='{}/{}'.format(mconfig['baseurl'], post['url']))
                feed_posts.append(item)

        # склеиваем всё и сохраняем в файл
        save_path = SAVE_TO + '/{}.rss'.format(feed['name'])
        with open(save_path, 'w', encoding='utf-8') as f:
            f.write('{}\n{}\n{}'.format(head, '\n'.join(feed_posts), foot))

        print(save_path, 'generated')


if __name__ == '__main__':
    run()
