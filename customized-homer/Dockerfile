FROM b4bz/homer:latest

RUN mkdir /tmp/build/
# Add context to /tmp/build/
COPY . /tmp/build

# this last command outputs the list of files added to the build context:
RUN find /tmp/build/

ADD assets /www/
EXPOSE 8080