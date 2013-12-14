from django.db import models

class Playlist(models.Model):
    name = models.CharField(max_length = 256)

class PlaylistSong(models.Model):
    models.ForeignKey(Playlist)
    path = models.CharField(max_length = 1024)
