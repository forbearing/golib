package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/forbearing/gst/ds/tree/trie"
	"github.com/spf13/cobra"
)

// RouteInfo represents information about a route
type RouteInfo struct {
	Path    string
	Methods []string
}

var routesCmd = &cobra.Command{
	Use:   "routes [filter]",
	Short: "print current route tree structure",
	Long: `Print the current route tree structure in a hierarchical format.
This command analyzes the registered routes and displays them as a tree structure.

Optional filter parameter can be used to show only routes matching the specified pattern.
For example: 'gg routes config/namespace' will show only routes under config/namespace.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var filter string
		if len(args) > 0 {
			filter = args[0]
		}
		routesRun(filter)
	},
}

func routesRun(filter string) {
	// Parse routes from generated router.go file
	routes, err := parseRoutesFromFile()
	if err != nil {
		fmt.Printf("Error parsing routes: %v\n", err)
		os.Exit(1)
	}

	if len(routes) == 0 {
		fmt.Println("No routes found. Please run 'gg gen' first to generate routes.")
		return
	}

	// Filter routes if filter is provided
	if filter != "" {
		routes = filterRoutes(routes, filter)
		if len(routes) == 0 {
			fmt.Printf("No routes found matching filter: %s\n", filter)
			return
		}
	}

	// Build and print route tree
	printRouteTree(routes)
}

// filterRoutes filters routes based on the given filter pattern
func filterRoutes(routes map[string][]string, filter string) map[string][]string {
	filteredRoutes := make(map[string][]string)

	// Normalize filter by removing leading/trailing slashes
	filter = strings.Trim(filter, "/")

	for path, methods := range routes {
		// Normalize path by removing leading slash
		normalizedPath := strings.TrimPrefix(path, "/")

		// Check if the path starts with the filter or contains the filter
		if strings.HasPrefix(normalizedPath, filter) || strings.Contains(normalizedPath, filter) {
			filteredRoutes[path] = methods
		}
	}

	return filteredRoutes
}

// parseRoutesFromFile parses the generated router.go file to extract route information
func parseRoutesFromFile() (map[string][]string, error) {
	// Check if router.go exists
	routerFile := filepath.Join(routerDir, "router.go")
	if !fileExists(routerFile) {
		return nil, fmt.Errorf("router file not found: %s. Please run 'gg gen' first", routerFile)
	}

	// Read and parse the router.go file
	content, err := os.ReadFile(routerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read router file: %w", err)
	}

	// Parse routes using regex patterns
	routes := make(map[string][]string)
	lines := strings.Split(string(content), "\n")

	// Pattern to match route registrations like: router.Register[*Type, *Type, *Type](router.Auth(), "path", &types.ControllerConfig[*Type]{...}, consts.Method)
	routePattern := regexp.MustCompile(`router\.Register\[.*?\]\([^,]+,\s*"([^"]+)",\s*&types\.ControllerConfig\[.*?\]\{.*?\},\s*consts\.(\w+)\)`)

	for _, line := range lines {
		matches := routePattern.FindStringSubmatch(strings.TrimSpace(line))
		if len(matches) == 3 {
			path := matches[1]
			method := matches[2]

			// Convert consts method to HTTP method
			httpMethod := convertConstToHTTPMethod(method)
			if httpMethod != "" {
				if routes[path] == nil {
					routes[path] = []string{}
				}
				routes[path] = append(routes[path], httpMethod)
			}
		}
	}

	return routes, nil
}

// convertConstToHTTPMethod converts consts method to HTTP method
func convertConstToHTTPMethod(constMethod string) string {
	switch constMethod {
	case "Create":
		return "POST"
	case "Delete":
		return "DELETE"
	case "Update":
		return "PUT"
	case "Patch":
		return "PATCH"
	case "List", "Get":
		return "GET"
	case "CreateMany":
		return "POST"
	case "DeleteMany":
		return "DELETE"
	case "UpdateMany":
		return "PUT"
	case "PatchMany":
		return "PATCH"
	case "Import", "Export":
		return "POST"
	default:
		return ""
	}
}

// printRouteTree builds and prints the route tree structure
func printRouteTree(routes map[string][]string) {
	// Create a trie to organize routes hierarchically
	trie, err := trie.New[string, *RouteInfo]()
	if err != nil {
		fmt.Printf("Error creating trie: %v\n", err)
		return
	}

	// Add all routes to the trie
	for path, methods := range routes {
		// Split path into segments
		segments := strings.Split(strings.Trim(path, "/"), "/")
		// Filter out empty segments
		var filteredSegments []string
		for _, segment := range segments {
			if segment != "" {
				filteredSegments = append(filteredSegments, segment)
			}
		}

		// Sort methods for consistent display
		sort.Strings(methods)

		// Add to trie
		routeInfo := &RouteInfo{
			Path:    path,
			Methods: methods,
		}
		trie.Put(filteredSegments, routeInfo)
	}

	// Print the tree
	fmt.Println("Route Tree Structure:")
	fmt.Println("=====================")
	if trie.IsEmpty() {
		fmt.Println("No routes found.")
		return
	}

	// Use custom formatter for better route display
	printTrieNode(trie.Root(), "", "", true)
}

// printTrieNode recursively prints the trie structure with improved formatting
func printTrieNode(node *trie.Node[string, *RouteInfo], prefix, childPrefix string, isRoot bool) {
	if node == nil {
		return
	}

	// Get all children and sort them
	children := node.Children()

	// Convert to slice for sorting
	type childPair struct {
		key  string
		node *trie.Node[string, *RouteInfo]
	}
	childList := make([]childPair, 0, len(children))
	for k, v := range children {
		childList = append(childList, childPair{k, v})
	}

	// Sort children by key
	sort.Slice(childList, func(i, j int) bool {
		return childList[i].key < childList[j].key
	})

	// Print children
	for i, child := range childList {
		isLast := i == len(childList)-1
		newPrefix := childPrefix

		// Check if this child has a route (is an endpoint)
		hasRoute := child.node.Value() != nil
		hasChildren := len(child.node.Children()) > 0

		if hasRoute && !hasChildren {
			// This is a terminal endpoint - show with route info
			route := child.node.Value()
			methodsStr := formatMethods(route.Methods)
			if isLast {
				fmt.Printf("%s└─ %s %s\n", childPrefix, child.key, methodsStr)
			} else {
				fmt.Printf("%s├─ %s %s\n", childPrefix, child.key, methodsStr)
			}
		} else {
			// This is a path segment - show as directory
			if isLast {
				fmt.Printf("%s└─ %s/\n", childPrefix, child.key)
				newPrefix += "   "
			} else {
				fmt.Printf("%s├─ %s/\n", childPrefix, child.key)
				newPrefix += "│  "
			}

			// If this node has a route but also has children, show the route info
			if hasRoute {
				route := child.node.Value()
				methodsStr := formatMethods(route.Methods)
				if isLast {
					fmt.Printf("%s   ● %s\n", childPrefix, methodsStr)
				} else {
					fmt.Printf("%s│  ● %s\n", childPrefix, methodsStr)
				}
			}

			// Recursively print children
			printTrieNode(child.node, "", newPrefix, false)
		}
	}
}

// formatMethods formats HTTP methods with colors for better visibility
func formatMethods(methods []string) string {
	var formatted []string
	for _, method := range methods {
		switch method {
		case "GET":
			formatted = append(formatted, "\033[32mGET\033[0m") // Green
		case "POST":
			formatted = append(formatted, "\033[34mPOST\033[0m") // Blue
		case "PUT":
			formatted = append(formatted, "\033[33mPUT\033[0m") // Yellow
		case "PATCH":
			formatted = append(formatted, "\033[35mPATCH\033[0m") // Magenta
		case "DELETE":
			formatted = append(formatted, "\033[31mDELETE\033[0m") // Red
		default:
			formatted = append(formatted, method)
		}
	}
	return "[" + strings.Join(formatted, ", ") + "]"
}
