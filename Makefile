build:
	go build -o thermon ./src

install:
	cp thermon /usr/local/bin

	if ! id -u thermon >/dev/null 2>&1; then \
		useradd --system --no-create-home --shell /usr/sbin/nologin thermon; \
	fi

	mkdir -p /usr/local/share/thermon
	chown thermon /usr/local/share/thermon

	cp thermon.service /etc/systemd/system
	systemctl daemon-reload
	systemctl enable thermon

clean:
	rm thermon