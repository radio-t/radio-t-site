#!/usr/bin/env python3

import re
import glob
from ilio import read, write
from datetime import datetime
import dateutil.parser
import frontmatter
from string import Template

# add image to content body and frontmatter
def extract_image(post_file, post):
    exp = r"^[\s\n]*!\[]\((.+)\)[\s\n]*"
    match = re.match(exp, post.content)
    image = match.group(1) if match else None

    if image and 'image' in post.metadata:
        if post.metadata['image'] != image:
            print(f'Different images in {post_file}')
            return post

    if not (image or 'image' in post.metadata):
        print(f'No images in {post_file}')
        return post

    image = image or post.metadata['image']

    # add image from post content
    post.content = re.sub(exp, '', post.content)
    post.content = Template('![]($image)\n\n$content').substitute(image=image, content=post.content)

    # add image to frontmatter
    post.metadata['image'] = image

    return post

def get_prep_link(number):
    prep_file = f'content/posts/prep-{number}.md'
    prep = frontmatter.load(prep_file)
    date = dateutil.parser.parse(prep.metadata['date'])
    permalink = date.strftime(f'/p/%Y/%m/%d/prep-{number}/')
    return permalink

def fix_prep_link(post_file, post):
    exp = r"http://new.radio-t.com/\d+/\d+/(\d+).html"
    match = re.search(exp, post.content)
    number = match.group(1) if match else None
    if number:
        post.content = re.sub(exp, get_prep_link(number), post.content)
        return post

def insert_final_newline(file):
    content = read(file)
    exp = r"\s*$"
    write(file, re.sub(exp, "\n", content))

def run():
    # Uncomment the sections you need to run
    # extract image
    # for post_file in glob.glob('content/posts/podcast-*'):
    #     post = frontmatter.load(post_file)
    #     post = extract_image(post_file, post)
    #     write(post_file, frontmatter.dumps(post))

    # fix prep links
    # for post_file in glob.glob('content/posts/podcast-*'):
    #     post = frontmatter.load(post_file)
    #     post = fix_prep_link(post_file, post)
    #     if post:
    #         write(post_file, frontmatter.dumps(post))

    for file in glob.glob('content/posts/*.md'):
        insert_final_newline(file)

if __name__ == '__main__':
    run()