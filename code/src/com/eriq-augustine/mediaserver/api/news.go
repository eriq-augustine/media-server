package api;

// Implements the "/news/" portion of the api.

import (
   "com/pollr/server/database"
   "com/pollr/server/log"
   "com/pollr/server/messages"
);

const DEFAULT_NEWS_COUNT = 5;

func getRecentNews(count int) (interface{}, error) {
   if (count <= 0) {
      count = DEFAULT_NEWS_COUNT;
   }

   newsArticles, err := database.GetRecentNews(count);
   if (err != nil) {
      log.ErrorE("Unable to get recent news articles", err);
      return "", err;
   }

   return messages.NewNewsArticles(newsArticles), nil;
}
