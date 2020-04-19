FROM ruby:2.6.5

RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update -qq && apt-get install -y build-essential nodejs yarn

RUN mkdir /app
WORKDIR /app
COPY . /app

RUN yarn install --check-files

# Define a default value for ARG RAILS_ENV.
ARG RAILS_ENV=development

# Install the dependencies for specific group.
RUN bundle install --with $RAILS_ENV
# Precompile assets for specific environment and place them under public/assets.
RUN bundle exec rake assets:precompile
