## Steps to clone dabdab.org

1. change remote and push to the new repo
2. change flytoml to point to the new app <projectname>
3. create the app on fly `flyctl apps create <projectname>`
4. create db, set 4gb ram `fly postgres create -n <projectname>db`
5. attach db to the app `flyctl postgres attach -a <projectname> <projectname>db`
6. Set secrets:

   ```
   flyctl secrets set SESSION_SALT=<random string>
   flyctl secrets set SITE_ROOT=https://<projectname>.com
   flyctl secrets set MJ_APIKEY_PUBLIC=<public key from mailjet>
   flyctl secrets set MJ_APIKEY_PRIVATE=<private key from mailjet>
   ```
7. Do first deploy `fly deploy`, make sure you can reach the app via <appname>.fly.dev
8. Create a cert for your custom domain `fly certs add <projectname>.com`
9. After it screams at you, add required A and AAAA records
10. You might need to run `fly certs check <projectname>.com` a couple of times, `fly certs list` should show your domain with the status `ready`.
11. You should be able to reach your app via custom domain at this point
12. Got to mailjet and add new domain
13. Add sender email address there
14. Add required txt record to validate domain
15. Add required txt records to add DKIM and SPF settings
16. Add postgres db env var to `cmd/web/.env` via `./env.pl > cmd/web/.env`, remove `sslMode=disable` and replace domain name with localhost
17. Move `pkg` folder back
18. Run `./generate.sh` to get model files
19. Work till you get the landing and login working
