var godan = angular.module('Godan', ['ngResource', 'ngAnimate', 'ngTouch', 'ui.bootstrap', 'ngRoute']);

var godan_api = "http://localhost:8000/";

godan.config(['$routeProvider', function($routeProvider) {
	$routeProvider.
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

godan.controller('NavCtrl', ['$scope', '$location', function NavCtrl($scope, $location) {
	$scope.isActive = function (viewLocation) { 
        return viewLocation === $location.path();
    };
}]);

godan.controller('QueriesCtrl', ["$scope", "$resource", "$uibModal", "$log", function QueriesCtrl($scope, $resource, $uibModal, $log) {

	$scope.animationsEnabled = true;

	$scope.submitQuery = function() {
		var query = godan_api + "ips";
		
		var params = {}
		if ($scope.inputPort) {
			params["port"] = $scope.inputPort;
		}
		if ($scope.inputService) {
			params["service"] = $scope.inputService;
		}
		if ($scope.inputRegexp) {
			params["regexp"] = $scope.inputRegexp;
		}

		if ($scope.inputIP) {
			query = query + "/" + $scope.inputIP
		}

		if (Object.keys(params).length > 0) {
			query += "?"+$.param(params);
		}

		$scope.resultsTable = $resource(query, {}, {}).query();
	}

	$scope.rowClicked = function(element) {
		var modalInstance = $uibModal.open({
			animation: $scope.animationsEnabled,
			templateUrl: 'modalTemplate',
			controller: 'ModalInstanceCtrl',
			size: 'lg',
			resolve: {
				element: function() {
					return element;
				}
			}
		});
		$scope.detailedContent = element.Content;
	};

}]);

godan.controller('ModalInstanceCtrl', function ($scope, $uibModalInstance, element) {
  $scope.closeit = function () {
    $uibModalInstance.close();
  };
	$scope.content = atob(element.Content); 
	$scope.elements = [element];
});

godan.controller('TasksCtrl', ["$scope", "$resource", "$uibModal", "$log", function TasksCtrl($scope, $resource, $uibModal, $log) {

	$scope.animationsEnabled = true;

	$scope.submitTask = function() {
		var url = godan_api + "tasks";
		var task = {
			"ips": $scope.taskIPs.split("\n"),
			"ports": $scope.taskPorts.split("\n"),
		}

		var taskResult = $resource(url, {}, {}).save(task);
		taskResult.$promise.then(function(result) {
			if (result.status == "success") {
				$uibModal.open({
					animation: $scope.animationsEnabled,
					templateUrl: 'taskModalTemplate',
					controller: 'ModalInstanceCtrl2',
					resolve: {
						title: function() { return "Success!"; },
						message: function() { return "Task successfully submitted!"; }
					}
				});
			} else {
				$uibModal.open({
					animation: $scope.animationsEnabled,
					templateUrl: 'taskModalTemplate',
					controller: 'ModalInstanceCtrl2',
					resolve: {
						title: function() { return "Fail!"; },
						message: function() { return "Task was not correctly submitted"; }
					}
				});
			}
		});

	}
}]);

godan.controller('ModalInstanceCtrl2', function ($scope, $uibModalInstance, title, message) {
  $scope.closeit = function () {
    $uibModalInstance.close();
  };
  $scope.title = title;
  $scope.message = message;
});

godan.controller('StatusCtrl', ["$scope", "$resource", "$interval", "$uibModal", "$log", function StatusCtrl($scope, $resource, $interval, $uibModal, $log) {
	var url = godan_api + "status";
	$scope.statusTable = $resource(url, {}, {}).query();

    refresh = $interval(function () {
	freshTable = $resource(url, {}, {}).query();
	freshTable.$promise.then(function(result) {
		console.log("Refresh!");
		$scope.statusTable = result;
		$scope.$applyAsync();
	});
  	}, 5000);

	$scope.$on('$locationChangeStart', function(){
		$interval.cancel(refresh);
	});

	$scope.rowClicked = function(element) {
		var modalInstance = $uibModal.open({
			animation: $scope.animationsEnabled,
			templateUrl: 'statusModalTemplate',
			controller: 'ModalInstanceCtrl3',
			size: 'lg',
			resolve: {
				element: function() {
					return element;
				},
				url: function() {
					return url;
				}
			}
		});
	};
}]);

godan.controller('ModalInstanceCtrl3', function ($scope, $resource, $uibModalInstance, $interval, element, url) {
  $scope.$on("modal.closing", function() {
	$interval.cancel(refresh);
  });

  $scope.closeit = function () {
    $uibModalInstance.close();
  };

  
  $scope.tasks = element.Tasks;
  $scope.info = element.Info;
  $scope.running = element.Running;
  refresh = $interval(function () {
	freshStatus = $resource(url, {}, {}).query();
	freshStatus.$promise.then(function(result) {
		console.log("ASD");
		$scope.tasks = result[0].Tasks;
		$scope.info = result[0].Info;
		$scope.running = result[0].Running;
		$scope.$applyAsync();
	});
  }, 2000);

  $scope.status = [
		{"name": "Pause", "value": "pause"},
		{"name": "Resume", "value": "resume"},
		{"name": "Soft shutdown", "value": "softshutdown"},
		{"name": "Hard shutdown", "value": "hardshutdown"}
  ];

  $scope.updateStatus = function() {
	  var data = {
			"target": element.Name,
			"command": $scope.newStatus,
	  }
	  var statusSet = $resource(url, {}, {}).save(JSON.stringify(data));
	  statusSet.$promise.then(function(result) {
			console.log("Status updated!")
	  });
  }
});

godan.controller('UsersCtrl', ["$scope", "$resource", "$uibModal", "$log", function UsersCtrl($scope, $resource, $uibModal, $log) {

}]);
