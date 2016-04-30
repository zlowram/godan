var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('QueriesCtrl', ["$scope", "$resource", "$uibModal", "$window", function QueriesCtrl($scope, $resource, $uibModal, $window) {
	if (!$window.sessionStorage.token) $window.location.href = '#/login';

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

		var queryResult = $resource(query, {}, {
			query: {
				method: 'GET',
				isArray: true,
				headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
			}
		}).query();
		queryResult.$promise.then(
			function(result) {
				if (result.length == 0) {
					$uibModal.open({
						animation: $scope.animationsEnabled,
						templateUrl: 'resultsModalTemplate',
						controller: 'ResultsModalInstanceCtrl',
						resolve: {
							title: function() { return "No results"; },
							message: function() { return "No entries matched the query"; }
						}
					});
				} else {
					$scope.resultsTable = result;
				}
			},
			function(result) {
				$uibModal.open({
					animation: $scope.animationsEnabled,
					templateUrl: 'resultsModalTemplate',
					controller: 'abcd',
					resolve: {
						title: function() { return "Error!"; },
						message: function() { return "There was an error sending the query"; }
					}
				});
			}
		);
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

angular.module('Godan').controller('ModalInstanceCtrl', function ($scope, $uibModalInstance, element) {
  $scope.closeit = function () {
    $uibModalInstance.close();
  };
	$scope.content = atob(element.Content); 
	$scope.elements = [element];
});

angular.module('Godan').controller('ResultsModalInstanceCtrl', function ($scope, $uibModalInstance, title, message) {
  $scope.closeit = function () {
    $uibModalInstance.close();
  };
  $scope.title = title;
  $scope.message = message;
});
