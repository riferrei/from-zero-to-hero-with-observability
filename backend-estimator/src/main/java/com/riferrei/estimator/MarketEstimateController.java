package com.riferrei.estimator;

import com.riferrei.estimator.domain.MarketEstimate;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class MarketEstimateController {

    @RequestMapping("/estimateValue")
    public MarketEstimate estimateValue(@RequestParam(name = "brand") String make,
                                   @RequestParam(name = "model") String model,
                                   @RequestParam(name = "year") int year) {
        MarketEstimate estimate = new MarketEstimate(make, model, year);
        estimate.calculateEstimate();
        return estimate;
    }
}
