.PHONY: test coverage coverage-html clean

# 运行所有测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 生成覆盖率报告
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -1

# 生成HTML覆盖率报告并在浏览器中打开
coverage-html:
	@echo "Generating HTML coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 清理生成的文件
clean:
	@echo "Cleaning up..."
	@rm -f coverage.out coverage.html
	@echo "Done"
