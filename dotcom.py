#!/usr/bin/env python

from datetime import datetime
from jinja2 import Environment
from jinja2 import FileSystemLoader
from jinja2 import select_autoescape
from marko.ext.gfm import gfm # github flavor markdown can do radio boxes
import frontmatter
import os
import re
import sys

# if environ var not set, place output files into website dir in root
OUTPUT_DIR = os.environ.get('DOTCOM_OUTPUT_DIR', '/var/www/dotcom')

def title_case(string):
    exceptions = ['a', 'an', 'the', 'in', 'to']
    words = re.split(' ', string)
    titled_list = list()
    titled_list.append(words[0].capitalize()) # first word is always capitalized
    for word in words[1:]:
        if word not in exceptions:
            word = word.capitalize()
        titled_list.append(word)
    titled_string = ' '.join(titled_list)
    return titled_string

def parse_markdown(md_file):
    post = dict()
    data = frontmatter.load(md_file)
    post['pub_date'] = data['date']
    post['content'] = gfm(data.content)
    return post

def render_template(md_filename, LAYOUT_DIR, TEMPLATE_DIR):
    input_filename = f'{LAYOUT_DIR}/{md_filename}'
    output_filename = f'{OUTPUT_DIR}/{md_filename}'.replace('.md', '.html')
    post = parse_markdown(input_filename)
    post['title'] = title_case(md_filename.replace('.md', ''))

    # this is only for getting the current year for the copyright portion in the footer of the website
    this_year = datetime.strftime(datetime.now(), '%Y')

    env = Environment(
        loader=FileSystemLoader(f'{TEMPLATE_DIR}'),
        autoescape=select_autoescape()
    )
    template = env.get_template(f'post.html')

    # plug everything into template and write to file
    output = template.render(post=post, this_year=this_year)
    with open(f'{output_filename}', 'w') as f:
        f.write(output)

def main():
    input_dir_items = os.listdir('/')
    LAYOUT_DIR = 'layout'
    TEMPLATE_DIR = 'templates'
    STATIC_DIR = 'static'

    layout_files = os.listdir(LAYOUT_DIR)

    for post in layout_files:
        if '.md' in post:
            render_template(post, LAYOUT_DIR, TEMPLATE_DIR)
        else:
            #TODO: parse recursively if directories are found
            pass
    # copy over static directory
    os.system(f'cp -r {STATIC_DIR} {OUTPUT_DIR}')

if __name__ == "__main__":
    main()
