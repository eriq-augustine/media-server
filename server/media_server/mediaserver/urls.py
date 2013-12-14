from django.conf.urls import patterns, url

from mediaserver import views

urlpatterns = patterns('',
   url(r'^$', views.index, name = 'index'),
   url(r'^home/?$', views.home, name = 'home'),
   url(r'^browse/(?P<path>.*)/?$', views.browse, name = 'browse'),
)
