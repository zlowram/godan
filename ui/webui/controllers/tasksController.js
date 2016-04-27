var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('TasksCtrl', ["$scope", "$resource", "$uibModal", "$window", function TasksCtrl($scope, $resource, $uibModal, $window) {

	$scope.animationsEnabled = true;

	$scope.submitTask = function() {
		var url = godan_api + "tasks";
		var task = {
			"ips": $scope.taskIPs.split("\n"),
			"ports": $scope.taskPorts.split("\n"),
		}

		var taskResult = $resource(url, {}, {
			save: {
				method: 'POST',
				headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
			}
		}).save(task);

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

angular.module('Godan').controller('ModalInstanceCtrl2', function ($scope, $uibModalInstance, title, message) {
  $scope.closeit = function () {
    $uibModalInstance.close();
  };
  $scope.title = title;
  $scope.message = message;
});
