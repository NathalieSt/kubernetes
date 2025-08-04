schema Endpoint:
    enabled?: bool
    name: str
    group?: str
    url: str
    method?: str
    # e.g. [STATUS] == 200
    conditions?: [str]
    interval?: str

