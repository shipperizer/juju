test_hooks() {
	if [ "$(skip 'test_hooks')" ]; then
		echo "==> TEST SKIPPED: hooks"
		return
	fi

	set_verbosity

	echo "==> Checking for dependencies"
	check_dependencies juju charmcraft

	file="${TEST_DIR}/test-hooks.log"

	bootstrap "test-hooks" "${file}"

	test_dispatching_script
	test_start_hook_fires_after_reboot

	destroy_controller "test-hooks"
}
