package com.riferrei.estimator.domain;

import java.util.Random;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@JsonIgnoreProperties(ignoreUnknown = true)
public class MarketEstimate {

    final static Logger logger = LoggerFactory.getLogger(MarketEstimate.class);
    final static Random random = new Random(System.currentTimeMillis());

    private int estimate = 0;
    private String brand;
    private String model;
    private int year;

    public MarketEstimate(String brand, String model, int year) {
        this.brand = brand;
        this.model = model;
        this.year = year;
    }

    public void calculateEstimate() {

        logger.info("Value estimation for brand: " + brand);

        int basePrice = 0;
        switch(brand) {
            case "Toyota":
                basePrice = 25000;
                break;
            case "Lexus":
                basePrice = 35000;
                break;
            case "Ford":
                basePrice = 20000;
                break;
            case "Nissan":
                basePrice = 20000;
                break;
            case "Tesla":
                basePrice = 60000;
                break;
            case "Ferrari":
                basePrice = specialPriceCalculation();
                break;
            default:
                basePrice = 30000;
        }

        int est = (int) ((((Math.random()) - .5) / 10 + 1) * basePrice);
        setEstimate(est);

    }

    private int specialPriceCalculation() {
        logger.debug("Calculating special price for exotic car...");
        logger.debug("Wait 5 seconds for market data to come in...");
        try {
            Thread.sleep(5000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        logger.debug("Use market data with a fixed base price");
        return random.nextInt(10) * 100000;
    }

    public int getEstimate() {
        return estimate;
    }

    public void setEstimate(int estimate) {
        this.estimate = estimate;
    }

    public String getBrand() {
        return brand;
    }

    public void setBrand(String brand) {
        this.brand = brand;
    }

    public String getModel() {
        return model;
    }

    public void setModel(String model) {
        this.model = model;
    }

    public int getYear() {
        return year;
    }

    public void setYear(int year) {
        this.year = year;
    }

}
