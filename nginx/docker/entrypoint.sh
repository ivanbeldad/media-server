echo "Executing entrypoint"

if [ -z "$(ls -A /etc/nginx/conf.d)" ]; then
  echo "Initializing nginx configuration"
  cp "/docker/default.conf" "/etc/nginx/conf.d/default.conf"
fi

echo "Starting nginx"

nginx -g "daemon off;"
