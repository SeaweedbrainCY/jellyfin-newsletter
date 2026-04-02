## Jellyfin test data

The folder contains the docker files used to spin up a test-dedicated jellyfin tenant to create original data for go-vcr. It should only be used to generate thoses cassettes, which shouldn't happen often. 

It only populates media fake files. Configuration of API keys and of the tenant need to be manually done. Some metadata such as folder (series folder/season folder) need to be manually fixed as well. 

The folder contains mounted volume used by the jellyfin container managed by go-testcontainer. 

It is a test-specific installation/configuration with fake media data.

All the secrets/token/configuration correspond to a temporary, test-only instance. There is no point of hidden them.

Only necessary data is committed. If fixture data change, it should be reflected in tests code

The movies/tvshows data (fake data with only a relevant name) are populated at runtime to manage the creation datetime used by Jellyfin to compute de DateCreated metadata.
