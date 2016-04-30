var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('StatusCtrl', ["$scope", "$resource", "$interval", "$uibModal", "$window", function StatusCtrl($scope, $resource, $interval, $uibModal, $window) {
	if (!$window.sessionStorage.token) $window.location.href = '#/login';

	$scope.animationsEnabled = true;

	var url = godan_api + "status";
	$scope.statusTable = $resource(url, {}, {
		query: {
			method: 'GET',
			isArray: true,
			headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
		}
	}).query();

    refresh = $interval(function () {
	freshTable = $resource(url, {}, {
		query:{
			method: 'GET',
			isArray: true,
			headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
		}
	}).query();

	freshTable.$promise.then(function(result) {
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

angular.module('Godan').controller('ModalInstanceCtrl3', function ($scope, $resource, $uibModalInstance, $interval, $window, element, url) {
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
	freshStatus = $resource(url, {}, {
		query: {
			method: 'GET',
			isArray: true,
			headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
		}
	}).query();

	freshStatus.$promise.then(function(result) {
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
	  var statusSet = $resource(url, {}, {
		  save: {
			  method: 'POST',
			  headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
		  }
	  }).save(JSON.stringify(data));
  }
});
