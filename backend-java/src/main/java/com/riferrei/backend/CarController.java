package com.riferrei.backend;

import java.util.Optional;

import com.riferrei.backend.domain.Car;
import com.riferrei.backend.domain.CarRepository;
import com.riferrei.backend.domain.MarketEstimate;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.util.UriComponentsBuilder;

@RestController
public class CarController {

    final static Logger logger = LoggerFactory.getLogger(CarController.class);

    @Value("${estimator.url}")
    private String estimatorUrl;

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

        UriComponentsBuilder builder = UriComponentsBuilder.fromUriString(estimatorUrl)
            .path("/estimateValue")
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
    public Car update(@PathVariable Long id, @RequestBody Car car) {

        logger.debug("In PUT by id");
        Optional<Car> optCar = carRepository.findById(id);
        Car updatedCar = optCar.get();

        if(car.getBrand() != null)
            updatedCar.setBrand(car.getBrand());
        if(car.getModel() != null)
            updatedCar.setModel(car.getModel());
        if(car.getColor() != null)
            updatedCar.setColor(car.getColor());
        if(car.getYear() != 0)
            updatedCar.setYear(car.getYear());
        if(car.getPrice() != 0)
            updatedCar.setPrice(car.getPrice());

        UriComponentsBuilder builder = UriComponentsBuilder.fromUriString(estimatorUrl)
            .path("/estimateValue")
            .queryParam("brand", updatedCar.getBrand())
            .queryParam("model", updatedCar.getModel())
            .queryParam("year", updatedCar.getYear());

        logger.debug(builder.build().toString());
        RestTemplate restTemplate = new RestTemplate();
        MarketEstimate carEstimate = restTemplate.getForObject(builder.build().toString(), MarketEstimate.class);
        logger.debug(carEstimate.toString());
        updatedCar.setMarketEstimate(carEstimate.getEstimate());

        carRepository.save(updatedCar);
        return updatedCar;

    }

    @RequestMapping(method=RequestMethod.DELETE, value="/api/cars/{id}")
    public Car delete(@PathVariable long id) {
        Optional<Car> optCar = carRepository.findById(id);
        Car car = optCar.get();
        carRepository.delete(car);
        return car;
    }

}
