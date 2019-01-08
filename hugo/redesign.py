#!/usr/bin/env python

import re
import glob
from ilio import read, write
from datetime import datetime
import frontmatter
from string import Template

# add image to content body and frontmatter
def extract_image(post_file, post):
    exp = r"^[\s\n]*!\[]\((.+)\)[\s\n]*"
    match = re.match(exp, post.content)
    image = match.group(1) if match else None

    if (image and ('image' in post.metadata)):
        if (post.metadata['image'] != image):
            print('Different images in ' + post_file)
            return post
    
    if (not (image or ('image' in post.metadata))):
        print('No images in ' + post_file)
        return post
    
    image = image if image else post.metadata['image']

    # add image from post content
    post.content = re.sub(exp, '', post.content)
    post.content = Template('![]($image)\n\n$content').substitute({'image': image, 'content': post.content})

    # add image to frontmatter
    post.metadata['image'] = image

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
