angular.module('Godan', ['ngResource', 'ngAnimate', 'ngTouch', 'ui.bootstrap', 'ngRoute']);

var godan_api = "http://localhost:8000/";

angular.module('Godan').config(['$httpProvider', function($httpProvider) {
	$httpProvider.defaults.withCredentials = false;	
}]);


angular.module('Godan').config(['$routeProvider', function($routeProvider) {
	$routeProvider.
		when('/login', {
			templateUrl: 'login.html',
			controller: 'LoginCtrl',
		}).
		when('/users', {
			templateUrl: 'users.html',
			controller: 'UsersCtrl',
		}).
		when('/tasks', {
			templateUrl: 'tasks.html',
			controller: 'TasksCtrl',
		}).
		when('/status', {
			templateUrl: 'status.html',
			controller: 'StatusCtrl',
		}).
		when('/queries', {
			templateUrl: 'ips.html',
			controller: 'QueriesCtrl',
		}).
		otherwise({
			redirectTo: '/queries'
		});
		
}]);

