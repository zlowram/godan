var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('StatusCtrl', ["$scope", "$resource", "$interval", "$uibModal", "$log", function StatusCtrl($scope, $resource, $interval, $uibModal, $log) {
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

angular.module('Godan').controller('ModalInstanceCtrl3', function ($scope, $resource, $uibModalInstance, $interval, element, url) {
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
