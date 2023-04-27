import argparse
import requests
import sys
import json

import os.path
from requests.exceptions import ConnectionError


api = "http://localhost:8080"


def allcontacts(auth_token):
    URL = api + '/contact/all' 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    try:
        resp = requests.get(
            url=URL,
            headers=HEADERS) 
    except ConnectionError:
        print("all contacts:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
    #logger.debug(" request : {}, {} ".format(URL, resp))      
    if resp.ok:
        result = resp.json()
        print(result)
    return

def authenticate(email, password):
    URL = api + '/authenticate' 
    HEADERS = {"Content-Type": "application/json", 'Accept': 'application/json'}
    DATA = {
        'email': email,
        'password': password
    }
    try:
        resp = requests.post(
            url=URL,
            data=json.dumps(DATA),
            headers=HEADERS) 
    except ConnectionError:
        print("all contacts:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
       
    
    result = resp.json()
    print(result)
    return

def newContact(first_name, last_name, full_name, address, mobile, email, password, skillid1, skillid2, auth_token):
    skills = []
    skills.append(skillid1)
    skills.append(skillid2)
    list_string = map(str, skills)
    skills_string = ' '.join(list_string)

    URL = api + '/contact/new' 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    DATA = {
        'first_name': first_name,
        'last_name': last_name,
        'full_name': full_name,
        'address': address,
        'email': email,
        'mobile': mobile,
        'password': password,
        'skillidsstring': skills_string,
    }
    try:
        resp = requests.put(
            url=URL,
            data=json.dumps(DATA),
            headers=HEADERS) 
    except ConnectionError:
        print("new contact:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
       
    if resp.ok:
        result = resp.json()
        print(result)
    return

def updateContact(id, first_name, last_name, full_name, address, mobile, email, skillid1, skillid2, auth_token):
    skills = []
    skills.append(skillid1)
    skills.append(skillid2)
    list_string = map(str, skills)
    skills_string = ' '.join(list_string)

    URL = api + "/contact/update/{}".format(id) 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    DATA = {
        'first_name': first_name,
        'last_name': last_name,
        'full_name': full_name,
        'address': address,
        'email': email,
        'mobile': mobile,
        'skillidsstring': skills_string,
    }
    try:
        resp = requests.patch(
            url=URL,
            data=json.dumps(DATA),
            headers=HEADERS) 
    except ConnectionError:
        print("new contact:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp.json())
       
    if resp.ok:
        result = resp.json()
        print(result)
    return


def deleteContact(id, auth_token):
    URL = api + "/contact/delete/{}".format(id) 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    try:
        resp = requests.delete(
            url=URL,
            headers=HEADERS) 
    except ConnectionError:
        print("delete contact:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
    #logger.debug(" request : {}, {} ".format(URL, resp))      
    if resp.ok:
        result = resp.json()
        print(result)
    return

def newSkill(name, level, auth_token):
    URL = api + '/skill/new' 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    DATA = {
        'name': name,
        'level': level
    }
    try:
        resp = requests.put(
            url=URL,
            data=json.dumps(DATA),
            headers=HEADERS) 
    except ConnectionError:
        print("all contacts:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
       
    if resp.ok:
        result = resp.json()
        print(result)
    return

def allSkills(auth_token):
    URL = api + '/skill/all' 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    try:
        resp = requests.get(
            url=URL,
            headers=HEADERS) 
    except ConnectionError:
        print("all skills:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
    #logger.debug(" request : {}, {} ".format(URL, resp))      
    if resp.ok:
        result = resp.json()
        print(result)
    return

def getSkill(id, auth_token):
    URL = api + "/skill/get/{}".format(id) 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    try:
        resp = requests.get(
            url=URL,
            headers=HEADERS) 
    except ConnectionError:
        print("all skills:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
    #logger.debug(" request : {}, {} ".format(URL, resp))      
    if resp.ok:
        result = resp.json()
        print(result)
    return

def getContact(id, auth_token):
    URL = api + "/contact/get/{}".format(id) 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    try:
        resp = requests.get(
            url=URL,
            headers=HEADERS) 
    except ConnectionError:
        print("get contact:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
    #logger.debug(" request : {}, {} ".format(URL, resp))      
    if resp.ok:
        result = resp.json()
        print(result)
    return

def updateSkill(id, name, level, auth_token):
    URL = api + "/skill/update/{}".format(id) 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    DATA = {
        'name': name,
        'level': level
    }
    try:
        resp = requests.patch(
            url=URL,
            data=json.dumps(DATA),
            headers=HEADERS) 
    except ConnectionError:
        print("update skill:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
    #logger.debug(" request : {}, {} ".format(URL, resp))      
    if resp.ok:
        result = resp.json()
        print(result)
    return

def deleteSkill(id, auth_token):
    URL = api + "/skill/delete/{}".format(id) 
    HEADERS = {"Content-Type": "application/json", 'Authorization': f'Bearer {auth_token}', 'Accept': 'application/json'}
    try:
        resp = requests.delete(
            url=URL,
            headers=HEADERS) 
    except ConnectionError:
        print("all skills:, OWT server connection failed, please check your network/firewall or report scicat server connection error to") 
        return 
    
    print(resp)
    #logger.debug(" request : {}, {} ".format(URL, resp))      
    if resp.ok:
        result = resp.json()
        print(result)
    return




if __name__ == "__main__":
    usage = "owt"
    parser = argparse.ArgumentParser(usage=usage)
    
    subparsers = parser.add_subparsers(dest="action")

    sub_auth = subparsers.add_parser("authenticate", help="authenticate user")
    sub_auth.add_argument('-e', '--email', dest='email', required=True, help="user email")
    sub_auth.add_argument('-p', '--password', dest='password', required=False, help="user password, prompted if not provided", default="")

    # contract operations
    sub_auth = subparsers.add_parser("contacts", help="Get all contacts")
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token") 

    sub_auth = subparsers.add_parser("newcontact", help="create new contact")
    sub_auth.add_argument('-f', '--first_name', dest='first_name', required=True, help="first name")
    sub_auth.add_argument('-l', '--last_name', dest='last_name', required=True, help="last name")
    sub_auth.add_argument('-u', '--full_name', dest='full_name', required=True, help="full name")
    sub_auth.add_argument('-a', '--address', dest='address', required=True, help="address")
    sub_auth.add_argument('-e', '--email', dest='email', required=True, help="user email")
    sub_auth.add_argument('-m', '--mobile', dest='mobile', required=True, help="mobile")
    sub_auth.add_argument('-p', '--password', dest='password', required=True, help="user password", default="")
    sub_auth.add_argument('-s1', '--skillid1', dest='skillid1', help="skillid 1", default=-1)
    sub_auth.add_argument('-s2', '--skillid2', dest='skillid2', help="skillid 2", default=-1)
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token")

    sub_auth = subparsers.add_parser("updatecontact", help="create new contact")
    sub_auth.add_argument('-i', '--id', dest='id', required=True, type=int, default=0, help="id")
    sub_auth.add_argument('-f', '--first_name', dest='first_name', required=True, help="first name")
    sub_auth.add_argument('-l', '--last_name', dest='last_name', required=True, help="last name")
    sub_auth.add_argument('-u', '--full_name', dest='full_name', required=True, help="full name")
    sub_auth.add_argument('-a', '--address', dest='address', required=True, help="address")
    sub_auth.add_argument('-e', '--email', dest='email', required=True, help="user email")
    sub_auth.add_argument('-m', '--mobile', dest='mobile', required=True, help="mobile")
    sub_auth.add_argument('-s1', '--skillid1', dest='skillid1', help="skillid 1", default=-1)
    sub_auth.add_argument('-s2', '--skillid2', dest='skillid2', help="skillid 2", default=-1)
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token")

    sub_auth = subparsers.add_parser("getcontact", help="get one contact")
    sub_auth.add_argument('-i', '--id', dest='id', required=True, type=int, default=0, help="id")
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token") \

    sub_auth = subparsers.add_parser("deletecontact", help="delete one contact")
    sub_auth.add_argument('-i', '--id', dest='id', required=True, type=int, default=0, help="id")
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token")

  
    # skill
    sub_auth = subparsers.add_parser("newskill", help="create new skill")
    sub_auth.add_argument('-n', '--name', dest='name', required=True, help="name")
    sub_auth.add_argument('-l', '--level', dest='level', required=True, type=int, default=0, help="level")
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token") 

    sub_auth = subparsers.add_parser("allskills", help="list all skills")
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token") 

    sub_auth = subparsers.add_parser("getskill", help="get one skill")
    sub_auth.add_argument('-i', '--id', dest='id', required=True, type=int, default=0, help="id")
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token") 

    sub_auth = subparsers.add_parser("deleteskill", help="delete one skill")
    sub_auth.add_argument('-i', '--id', dest='id', required=True, type=int, default=0, help="id")
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token")

    sub_auth = subparsers.add_parser("updateskill", help="update one skill")
    sub_auth.add_argument('-i', '--id', dest='id', required=True, type=int, default=0, help="id")
    sub_auth.add_argument('-n', '--name', dest='name', required=True, help="name")
    sub_auth.add_argument('-l', '--level', dest='level', required=True, type=int, default=0, help="level")
    sub_auth.add_argument('-t', '--token', dest='token', required=True, help="user token") 
    
    args = parser.parse_args()
    
    if args.action is None:
        parser.print_help()
        sys.exit()
    ## modified for issue 28
    res = None
    if args.action == "contacts":
        allcontacts(args.token) 
    elif args.action == "authenticate":
        authenticate(args.email, args.password)      
    elif args.action == "newcontact":
        newContact(args.first_name, args.last_name, args.full_name, args.address, args.mobile, args.email, args.password, args.skillid1, args.skillid2, args.token)          
    elif args.action == "newskill":
        newSkill(args.name, args.level, args.token)              
    elif args.action == "getskill":
        getSkill(args.id, args.token)                  
    elif args.action == "updateskill":
        updateSkill(args.id, args.name, args.level, args.token)                      
    elif args.action == "deleteskill":
        deleteSkill(args.id, args.token)                      
    elif args.action == "deletecontact":
        deleteContact(args.id, args.token)  
    elif args.action == "updatecontact":
        updateContact(args.id, args.first_name, args.last_name, args.full_name, args.address, args.mobile, args.email, args.skillid1, args.skillid2, args.token)                
    elif args.action == "getcontact":
        getContact(args.id, args.token)    
    