#!/usr/bin/env python3

"""
Скрипт для генерации rss-файлов

pip install pytoml mistune
"""

import glob
import subprocess
import argparse

from mistune import create_markdown
from mistune.renderers import HTMLRenderer

import pytoml as toml
from datetime import datetime

POSTS_DIR = './content/posts'
SAVE_TO = '/srv/hugo/public'
# SAVE_TO = '/tmp'
DATA_RSS = './data/rss'
FEEDS = [
    {'name': 'podcast', 'title': 'Радио-Т',
     'image': 'https://radio-t.com/images/covers/cover.png', 'count': 20, 'size': True},
    {'name': 'podcast-archives', 'title': 'Радио-Т Архивы',
     'image': 'https://radio-t.com/images/covers/cover-archive.png', 'count': 1000, 'size': False},
    {'name': 'podcast-archives-short', 'title': 'Радио-Т Архивы',
     'image': 'https://radio-t.com/images/covers/cover-archive.png', 'count': 25, 'size': False},
]


def parse_args():
    parser = argparse.ArgumentParser(description="Generate RSS files")
    parser.add_argument("--save-to", dest="save_to", default="/srv/hugo/public", help="Override save directory")
    return parser.parse_args()


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
    renderer = HTMLRenderer(escape=False)
    markdown = create_markdown(renderer=renderer)

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

                # добавляем строчку "Темы" в описание эпизода перед списком тем, чтобы верстка списка не ехала в Apple Podcasts,
                # см. примеры в https://github.com/radio-t/radio-t-site/pull/128
                # при этом summary надо оставить как есть, т.к. оно показывается в одну строчку и там проблемы с версткой нет
                post_description_html = markdown(post['data'])
                rss_description_html = (post_description_html
                    .replace('<ul>', '<p><em>Темы</em><ul>', 1)
                    .replace('</ul>', '</ul></p>', 1))

                item = body.format(title=post['config']['title'], 
                                   description=rss_description_html,
                                   summary=post_description_html,
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
    args = parse_args()
    SAVE_TO = args.save_to  # Override SAVE_TO if --save-to is provided
    run()
