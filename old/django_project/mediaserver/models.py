from django.db import models

class Playlist(models.Model):
    name = models.CharField(max_length = 256)

class PlaylistSong(models.Model):
    playlist = models.ForeignKey(Playlist)
    path = models.CharField(max_length = 1024)

class EncodeQueue(models.Model):
    queue_time = models.DateTimeField(auto_now_add = True)
    src = models.CharField(max_length = 2048)
    hash = models.CharField(max_length = 32)

class Cache(models.Model):
    cache_time = models.DateTimeField(auto_now_add = True)
    src = models.CharField(max_length = 2048)
    hash = models.CharField(max_length = 32)
    urlpath = models.CharField(max_length = 2048)
    hit_count = models.IntegerField(default = 1)
    last_access = models.DateTimeField(auto_now = True)
    bytes = models.IntegerField()
