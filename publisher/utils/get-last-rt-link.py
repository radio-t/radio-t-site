#!/usr/bin/python
# -*- coding: utf-8 -*-

import os, string, time, smtplib, shutil, stat, urllib, glob


if __name__ == "__main__":
    line = os.popen("curl https://radio-t.com/ | grep podcast- | head -n1").readline()
    link = "https://radio-t.com" + line.split("\"")[3]
    print link
