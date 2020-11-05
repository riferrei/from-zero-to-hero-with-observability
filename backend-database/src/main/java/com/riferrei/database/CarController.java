package com.riferrei.database;

import java.util.Optional;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.util.UriComponentsBuilder;

import com.riferrei.database.domain.CarRepository;
import com.riferrei.database.domain.Car;
import com.riferrei.database.domain.MarketEstimate;

@RestController
public class CarController {

    final static Logger logger = LoggerFactory.getLogger(CarController.class);

    @Value("${estimator.uri}")
    private String estimatorUri;

    @Autowired
    CarRepository carRepository;

    @Autowired
    public CarController(CarRepository repo) {
        carRepository = repo;
    }

    @RequestMapping(method= RequestMethod.GET, value="/api/cars")
    public Iterable<Car> Car() {
        logger.debug("In GET All");
        return carRepository.findAll();
    }

    @RequestMapping(method=RequestMethod.POST, value="/api/cars")
    public Car save(@RequestBody Car car) {

        logger.debug("In POST add");

        UriComponentsBuilder builder = UriComponentsBuilder.fromUriString(estimatorUri)
            .queryParam("brand", car.getBrand())
            .queryParam("model", car.getModel())
            .queryParam("year", car.getYear());

        logger.debug(builder.build().toString());
        RestTemplate restTemplate = new RestTemplate();
        MarketEstimate carEstimate = restTemplate.getForObject(builder.build().toString(), MarketEstimate.class);
        logger.debug(carEstimate.toString());

        car.setMarketEstimate(carEstimate.getEstimate());
        carRepository.save(car);
        return car;

    }

    @RequestMapping(method=RequestMethod.GET, value="/api/cars/{id}")
    public Optional<Car> show(@PathVariable Long id) {
        logger.debug("In GET by id");
        return carRepository.findById(id);
    }

    @RequestMapping(method=RequestMethod.PUT, value="/api/cars/{id}")
    public Car update(@PathVariable Long id, @RequestBody Car Car) {

        logger.debug("In PUT by id");

        Optional<Car> optCar = carRepository.findById(id);
        Car car = optCar.get();
        if(Car.getBrand() != null)
            car.setBrand(Car.getBrand());
        if(Car.getModel() != null)
            car.setModel(Car.getModel());
        if(Car.getColor() != null)
            car.setColor(Car.getColor());
        if(Car.getYear() != 0)
            car.setYear(Car.getYear());
        if(Car.getPrice() != 0)
            car.setPrice(Car.getPrice());

        UriComponentsBuilder builder = UriComponentsBuilder.fromUriString(estimatorUri)
                .queryParam("brand", car.getBrand())
                .queryParam("model", car.getModel())
                .queryParam("year", car.getYear());

        logger.debug(builder.build().toString());
        RestTemplate restTemplate = new RestTemplate();
        MarketEstimate carEstimate = restTemplate.getForObject(builder.build().toString(), MarketEstimate.class);
        logger.debug(carEstimate.toString());
        car.setMarketEstimate(carEstimate.getEstimate());

        carRepository.save(car);
        return car;

    }

    @RequestMapping(method=RequestMethod.DELETE, value="/api/cars/{id}")
    public Car delete(@PathVariable long id) {
        Optional<Car> optCar = carRepository.findById(id);
        Car car = optCar.get();
        carRepository.delete(car);
        return car;
    }

}
