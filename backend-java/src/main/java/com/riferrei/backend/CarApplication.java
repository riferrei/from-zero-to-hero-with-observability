package com.riferrei.backend;

import com.riferrei.backend.domain.Car;
import com.riferrei.backend.domain.CarRepository;
import com.riferrei.backend.domain.Owner;
import com.riferrei.backend.domain.OwnerRepository;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

@SpringBootApplication
public class CarApplication {

	private static final Logger logger = LoggerFactory.getLogger(CarApplication.class);

	@Autowired
	private CarRepository repository;

	@Autowired
	private OwnerRepository orepository;

	public static void main(String[] args) {
		SpringApplication.run(CarApplication.class, args);
		logger.info("Car Application with Spring Boot");
	}

	@Bean
	CommandLineRunner runner() {

		return args -> {

			Owner owner1 = new Owner("John", "Johnson");
			Owner owner2 = new Owner("Mary", "Robinson");
			orepository.save(owner1);
			orepository.save(owner2);

			Car car = new Car("Ford", "Mustang", "Red", "ADF-1121", 2017, 59000, 58321, owner1);
			repository.save(car);
			car = new Car("Nissan", "Leaf", "White", "SSJ-3002", 2014, 29000, 31998, owner2);
			repository.save(car);
			car = new Car("Toyota", "Prius", "Silver", "KKO-0212", 2018, 39000, 41556, owner2);
			repository.save(car);
			
		};
		
	}

}
