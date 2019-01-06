#!/usr/bin/env python

import re
import glob
from ilio import read, write
from datetime import datetime
import frontmatter

def extract_image(post_file, post):
    # remove image from content
    # add image to frontmatter
    exp = r"^[\s\n]*!\[]\((.+)\)[\s\n]*"
    match = re.match(exp, post.content)
    if match:
        if ('image' in post.metadata):
            if (post.metadata['image'] != match.group(1)):
                print('Different images in ' + post_file)
        else:
            post.metadata['image'] = match.group(1)
        post.content = re.sub(exp, '', post.content)

    return post

def run():
    for post_file in glob.glob('content/posts/podcast-*'):
    # for post_file in ['content/posts/podcast-606.md']:
        post = frontmatter.load(post_file)
        post = extract_image(post_file, post)
        write(post_file, frontmatter.dumps(post))

        # break


if __name__ == '__main__':
    run()
