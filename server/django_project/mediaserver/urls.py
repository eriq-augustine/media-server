from django.conf.urls import patterns, url

from mediaserver import views

urlpatterns = patterns('',
   url(r'^$', views.index, name = 'index'),
   url(r'^home/?$', views.home, name = 'home'),
   url(r'^browse/?$', views.browse, name = 'browse'),
   url(r'^browse/(?P<urlpath>.*)/?$', views.browse, name = 'browse'),
   url(r'^view/(?P<urlpath>.*)/?$', views.view, name = 'view'),
   url(r'^raw/(?P<urlpath>.*)/?$', views.raw, name = 'raw'),
   url(r'^cache/(?P<urlpath>.*)/?$', views.fetch_cache, name = 'fetch_cache'),
   url(r'^gallery/(?P<urlpath>.*)/?$', views.gallery, name = 'gallery'),
)
