# -*- coding: utf-8 -*-
"""
Created on Fri Jun 09 10:15:08 2017

@author: williamfl
"""
import requests
import simplejson
import os
import sys
from time import sleep
from datetime import datetime, timedelta
from dateutil.tz import tzlocal
from glob import glob
import configparser
import argparse

os.chdir(os.path.dirname(__file__))
print(os.getcwd())
config = configparser.ConfigParser()
config.read('config.ini')


parser = argparse.ArgumentParser(description='Scheduler Master')
parser.add_argument('-s', dest='start', action='store_true',help='Start an instance now')
parser.add_argument('-d', dest='debug', action='store_true',help='Set logging to Debug level and output to console')
parser.add_argument('--id',type=str,help='id for job')
parser.add_argument('--single-delete',dest='sd',action='store_true')
parser.add_argument('--delete-all-from-cache',dest='dafc',action='store_true')
parser.add_argument('--list-contents',action='store_true')

pKeys = [u'LastRunAt',
 u'Domain',
 u'Name',
 u'Schedule',
 u'ScheduleTime',
 u'epsilon',
 u'Next Run At',
 u'Active',
 u'Command',
 u'TimesToRepeat',
 u'Owner',
 u'Application',
 u'SubDomain',
 u'Type',
 u'ID',
 u'InternalType']
API_URL = config.get('info',"API_URL")
def next_weekday(d, weekday):
    days_ahead = weekday - d.weekday()
    if days_ahead <= 0: # Target day already happened this week
        days_ahead += 7
    return d + timedelta(days_ahead)

class DataModel:
    def __init__(self,name):
        if name == "":
            raise Exception("Name cannot be left blank")
        self.name = name
        self.data = {"name":name}
    def setData(self,key,value):
        if not key in pKeys:
            return False
        self.data[key] = value
        return True
class Scheduler:
    cachelocation = config.get('info','CACHE_LOCATION')
    cache = None
    def __init__(self,url=API_URL):
        self.API_URL = url
        try:
            self.load_cache()
        except:
            self.cache = {}

    def add(self,data,replace=False):
        "Add Job that is to be executed"
        r = requests.post(self.API_URL, data=simplejson.dumps({'Reqjob':data}),
                    headers={'TYPE':'UNIQUE'})
        if r.status_code != 201:
            if r.status_code > 300:
                raise Exception(r.json()['error'])
            if not replace:
                return None
            else:
                s = r.json()
                id = s['error'].split(';')[1]
                self.replace(data,id)
                return
        elif r.status_code == 201:
            self.add_to_cache(r,data)
        return r.json()['Id']
    
    def replace(self,data,id):
        r = requests.put(self.API_URL+id,data=simplejson.dumps({'Reqjob':data}))
        if r.status_code != 200:
            raise Exception(r.text)
    def load_cache(self):
        if not os.path.exists(os.path.dirname(self.cachelocation)):
            try:
                os.makedirs(os.path.dirname(self.cachelocation))
            except OSError as exc: # Guard against race condition
                if exc.errno != errno.EEXIST:
                    raise
        with open(self.cachelocation,'r') as f:
            self.cache = simplejson.load(f)
    def save_cache(self):
        with open(self.cachelocation,'w') as f:
            simplejson.dump(self.cache,f,indent=4,sort_keys=True)
    
    def update_cache(self):
        r = requests.get(self.API_URL)
        t = r.json()['Results']
        v = {}
        for x in t:
            v[x['ID']] = x
        self.cache = v
        self.save_cache()
    def add_to_cache(self,req,data):
        self.cache[req.json()['Id']] = data
        self.save_cache()
    def get(self,id):
        return requests.get(self.API_URL + id).json()
    def delete(self,id):
        r = requests.post(self.API_URL + "remove/" + id)
        if r.status_code != 200:
            raise Exception(r.text)
    def delete_from_cache(self):
        for key, value in self.cache.iteritems():
            self.delete(key)
        self.cache = {}
        self.save_cache()
    def single_delete_from_cache(self,id):
        keys = []
        for key in self.cache.keys():
            if key.startswith(id):
                keys.append(key)
        if len(keys) == 1:
            del self.cache[keys[0]]
            self.save_cache()
            self.delete(keys[0])
        else:
            print("Cannot delete id %s, found %d instances" %(id,len(keys)))
    
    def get_list(self):
        self.update_cache()
        for key in self.cache.keys():
            print(' '.join([key," - ",self.cache[key]['Name']]))

    

if __name__ == "__main__":
    args = parser.parse_args()
    s = Scheduler()
    if args.sd:
        s.single_delete_from_cache(args.id)
    if args.list_contents:
        s.get_list()
    #i = raw_input('Empty cache')
    #if i == 'y':
    #s.delete_from_cache()
