package com.distributedaddition.serviceb.service;
import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import lombok.AllArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/numbers")
public class FileController {
    private final FileStorage fileStorage;
    private final Counter addNumberCounter;
    private final Counter getLastNumberCounter;

    public FileController(FileStorage fileStorage, MeterRegistry meterRegistry) {
        this.fileStorage = fileStorage;
        this.addNumberCounter = Counter.builder("numbers_add_requests_total")
                .description("Total requests to /numbers/add")
                .register(meterRegistry);
        this.getLastNumberCounter = Counter.builder("numbers_last_requests_total")
                .description("Total requests to /numbers/last")
                .register(meterRegistry);
    }
    @PostMapping("/add")
    public ResponseEntity<String> addNumber(@RequestParam int value){
        fileStorage.addNumber(value);
        addNumberCounter.increment();
        return ResponseEntity.ok("Number added successfully");
    }

    @GetMapping("/last")
    public ResponseEntity<Integer> getLastNumber(){
        int lastNumber = fileStorage.getLastNumber();
        getLastNumberCounter.increment();
        return ResponseEntity.ok(lastNumber);
    }

}
